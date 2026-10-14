/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import DefaultNumberOfBatches
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

class ParallelTestsFeatureTests {

    private val originalBatches = ServiceNumberOfBatches

    @After
    fun restoreServiceBatches() {
        ServiceNumberOfBatches = originalBatches
    }

    private fun collectBuilds(root: Project): List<BuildType> =
        root.buildTypes + root.subProjects.flatMap { collectBuilds(it) }

    private fun isAcceptanceBuild(build: BuildType): Boolean =
        build.name.endsWith(" - Acceptance Tests")

    @Test
    fun parallelTestsFeatureEnabledByDefault() {
        ServiceNumberOfBatches = mapOf()
        val root = googleCloudRootProject(testContextParameters())
        val acceptanceBuilds = collectBuilds(root).filter { isAcceptanceBuild(it) }
        assertTrue("There should be acceptance-test builds", acceptanceBuilds.isNotEmpty())
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
        assertTrue("There should be VCR recording builds", builds.any { it.name.startsWith("VCR Recording - ") })
        builds.filter { !isAcceptanceBuild(it) }.forEach { bt ->
            assertTrue("Non-test build '${bt.id}' should not enable Parallel Tests", bt.features.items.filterIsInstance<ParallelTestsFeature>().isEmpty())
        }
    }

    @Test
    fun everyTestRunnerShouldRunOnlyItsBatch() {
        collectBuilds(googleCloudRootProject(testContextParameters())).filter { isAcceptanceBuild(it) }.forEach { bt ->
            val scripts = bt.steps.items.flatMap { step -> step.params.filter { it.name == "script.content" }.map { it.value } }
            assertTrue("Build '${bt.id}' should run tests through scripts/teamcitytestscripts with -parallel-tests", scripts.any {
                it.contains("teamcity-test-runner\" -mode acceptance -parallel-tests ") || it.contains("teamcity-diff-test -parallel-tests ")
            })
        }
    }
}
