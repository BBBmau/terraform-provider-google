/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import DefaultNumberOfBatches
import builds.parallelTestsBatchScript
import builds.parallelTestsListScript
import generated.ServiceNumberOfBatches
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.BuildTypeSettings
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.buildFeatures.ParallelTestsFeature
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import projects.googleCloudRootProject
import java.nio.file.Files

class ParallelTestsFeatureTests {

    private val originalBatches = ServiceNumberOfBatches

    @After
    fun restoreServiceBatches() {
        ServiceNumberOfBatches = originalBatches
    }

    private fun collectBuilds(root: Project): List<BuildType> =
        root.buildTypes + root.subProjects.flatMap { collectBuilds(it) }

    private fun isAcceptanceBuild(build: BuildType): Boolean =
        build.name.endsWith(" - Acceptance Tests") || build.name.startsWith("VCR Recording - ")

    @Test
    fun parallelTestsFeatureEnabledByDefault() {
        ServiceNumberOfBatches = mapOf()
        val root = googleCloudRootProject(testContextParameters())
        val acceptanceBuilds = collectBuilds(root).filter { isAcceptanceBuild(it) }
        assertTrue("There should be acceptance-test builds", acceptanceBuilds.isNotEmpty())
        assertTrue("There should be VCR recording builds", acceptanceBuilds.any { it.name.startsWith("VCR Recording - ") })
        acceptanceBuilds.forEach { bt ->
            val feature = bt.features.items.filterIsInstance<ParallelTestsFeature>().singleOrNull()
            assertTrue("Build '${bt.id}' should enable the Parallel Tests feature by default", feature != null)
            assertEquals("Build '${bt.id}' should use the default number of batches", DefaultNumberOfBatches, feature!!.numberOfBatches)
        }
    }

    @Test
    fun batchOverrideCustomizesNumberOfBatches() {
        ServiceNumberOfBatches = mapOf("compute" to 8)
        val root = googleCloudRootProject(testContextParameters())

        val builds = collectBuilds(root).filter { isAcceptanceBuild(it) }
        assertTrue("There should be compute package builds", builds.any { it.params.findRawParam("PACKAGE_PATH")?.value?.endsWith("/compute") == true })
        builds.forEach { bt ->
            val expected = if (bt.params.findRawParam("PACKAGE_PATH")?.value?.endsWith("/compute") == true) 8 else DefaultNumberOfBatches
            assertEquals("Build '${bt.id}' should honor its batch count", expected, bt.features.items.filterIsInstance<ParallelTestsFeature>().single().numberOfBatches)
        }
    }

    @Test
    fun singleBatchOverrideDisablesFeatureForThatService() {
        ServiceNumberOfBatches = mapOf("compute" to 1)
        val root = googleCloudRootProject(testContextParameters())

        val builds = collectBuilds(root).filter { isAcceptanceBuild(it) }
        assertTrue("There should be compute package builds", builds.any { it.params.findRawParam("PACKAGE_PATH")?.value?.endsWith("/compute") == true })
        builds.forEach { bt ->
            val expected = if (bt.params.findRawParam("PACKAGE_PATH")?.value?.endsWith("/compute") == true) 0 else 1
            assertEquals("Only compute builds should disable batching", expected, bt.features.items.filterIsInstance<ParallelTestsFeature>().size)
        }
    }

    @Test
    fun nonTestBuildsShouldNotBeBatched() {
        val builds = collectBuilds(googleCloudRootProject(testContextParameters()))
        assertTrue("There should be composites", builds.any { it.type == BuildTypeSettings.Type.COMPOSITE })
        builds.filter { !isAcceptanceBuild(it) }.forEach { bt ->
            assertTrue("Non-test build '${bt.id}' should not enable Parallel Tests", bt.features.items.filterIsInstance<ParallelTestsFeature>().isEmpty())
        }
    }

    @Test
    fun everyTestRunnerShouldFilterBatches() {
        collectBuilds(googleCloudRootProject(testContextParameters())).filter { isAcceptanceBuild(it) }.forEach { bt ->
            assertTrue("Build '${bt.id}' should consume the batch exclusions", bt.steps.items.any {
                it.params.any { param -> param.name == "script.content" && param.value.contains(parallelTestsBatchScript()) }
            })
        }
    }

    private fun runBatchScripts(exclusions: String?, script: String): Pair<Int, String> {
        val directory = Files.createTempDirectory("teamcity-parallel-tests").toFile()
        try {
            val excludesFile = directory.resolve("excludes.txt")
            if (exclusions != null) excludesFile.writeText(exclusions)
            val path = if (exclusions == null) "" else excludesFile.absolutePath
            val process = ProcessBuilder("bash", "-c",
                parallelTestsBatchScript().replace("%system.teamcity.build.parallelTests.excludesFile%", path) + "\n" + script
            ).directory(directory).redirectErrorStream(true).start()
            val output = process.inputStream.bufferedReader().use { it.readText() }
            return Pair(process.waitFor(), output)
        } finally {
            directory.deleteRecursively()
        }
    }

    @Test
    fun customRunnersShouldExcludeOnlyExactTestNames() {
        val script = parallelTestsListScript().replace(
            "./test-binary -test.list=\"%TEST_PREFIX%\"",
            "printf 'TestAccOne\\nTestAccOneMore\\nTestAccTwo\\n'"
        ) + "\ncat batch-tests.txt"
        val (status, output) = runBatchScripts("# metadata\n\nTestAccOne\r\n", script)
        assertEquals(output, 0, status)
        assertEquals("TestAccOneMore\nTestAccTwo\n", output)
    }

    @Test
    fun emptyExclusionsShouldKeepAllTests() {
        listOf(null, "", "# metadata\n \n").forEach { exclusions ->
            val script = parallelTestsListScript().replace(
                "./test-binary -test.list=\"%TEST_PREFIX%\"",
                "printf 'TestAccOne\\nTestAccTwo\\n'"
            ) + "\ncat batch-tests.txt"
            val (status, output) = runBatchScripts(exclusions, script)
            assertEquals(output, 0, status)
            assertEquals("TestAccOne\nTestAccTwo\n", output)
        }
    }

    @Test
    fun completelyExcludedBatchShouldHaveNoTests() {
        val script = parallelTestsListScript().replace(
            "./test-binary -test.list=\"%TEST_PREFIX%\"",
            "printf 'TestAccOne\\nTestAccTwo\\n'"
        ) + "\ncat batch-tests.txt"
        val (status, output) = runBatchScripts("TestAccOne\nTestAccTwo\n", script)
        assertEquals(output, 0, status)
        assertEquals("", output)
    }

    @Test
    fun nativeGoRunnerShouldPreserveUserFlagsAndEscapeExclusions() {
        listOf("-skip=TestAccManual", "-skip TestAccManual", "-test.skip=TestAccManual").forEach { skip ->
            val (status, output) = runBatchScripts("TestAccOne.More\n", """
                go() { printf '%s\n' "${'$'}@"; }
                runBatchGoTests ./google/services/... -run=TestAcc $skip
            """.trimIndent())
            assertEquals(output, 0, status)
            assertEquals("test\n./google/services/...\n-run=TestAcc\n-skip=^(TestAccOne\\.More)${'$'}|(TestAccManual)\n", output)
        }
    }

    @Test
    fun subtestSkipFiltersShouldNotSilentlyBreakBatchExclusions() {
        val script = """
            go() { printf '%s\n' "${'$'}@"; }
            runBatchGoTests ./google/services/... -skip=TestAccOne/subtest
        """.trimIndent()
        val (status, output) = runBatchScripts("TestAccTwo\n", script)
        assertEquals(output, 1, status)
        assertTrue(output, output.contains("Cannot combine a subtest -skip filter"))
        val (singleBatchStatus, singleBatchOutput) = runBatchScripts(null, script)
        assertEquals(singleBatchOutput, 0, singleBatchStatus)
        assertEquals("test\n./google/services/...\n-skip=TestAccOne/subtest\n", singleBatchOutput)
    }

    @Test
    fun missingExclusionsShouldFailInsteadOfRunningAllTests() {
        val (status, output) = runBatchScripts(null, """
            EXCLUDES_FILE=/does/not/exist
            ${parallelTestsBatchScript().replace("EXCLUDES_FILE=\"%system.teamcity.build.parallelTests.excludesFile%\"", "")}
        """.trimIndent())
        assertEquals(output, 1, status)
        assertTrue(output, output.contains("Cannot read Parallel Tests exclusions"))
    }
}
