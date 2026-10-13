/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import DefaultNumberOfBatches
import ServiceSweeperName
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

    // collectPackageBuilds returns the per-service acceptance test builds in the GA and Beta
    // Nightly Tests projects, excluding composite and sweeper builds.
    private fun collectPackageBuilds(root: Project): List<BuildType> {
        val gaNightly = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)
        val betaNightly = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)
        return (gaNightly.buildTypes + betaNightly.buildTypes).filter { bt ->
            bt.type != BuildTypeSettings.Type.COMPOSITE && !bt.name.startsWith(ServiceSweeperName)
        }
    }

    @Test
    fun parallelTestsFeatureEnabledByDefault() {
        ServiceNumberOfBatches = mapOf()
        val root = googleCloudRootProject(testContextParameters())
        collectPackageBuilds(root).forEach { bt ->
            val feature = bt.features.items.filterIsInstance<ParallelTestsFeature>().singleOrNull()
            assertTrue("Build '${bt.id}' should enable the Parallel Tests feature by default", feature != null)
            assertEquals("Build '${bt.id}' should use the default number of batches", DefaultNumberOfBatches, feature!!.numberOfBatches)
        }
    }

    @Test
    fun batchOverrideCustomizesNumberOfBatches() {
        ServiceNumberOfBatches = mapOf("compute" to 8)
        val root = googleCloudRootProject(testContextParameters())

        val overriddenBuilds = collectPackageBuilds(root).filter { bt ->
            val feature = bt.features.items.filterIsInstance<ParallelTestsFeature>().singleOrNull()
            feature?.numberOfBatches == 8
        }

        assertTrue("At least one build should use the overridden number of batches", overriddenBuilds.isNotEmpty())
    }

    @Test
    fun singleBatchOverrideDisablesFeatureForThatService() {
        ServiceNumberOfBatches = mapOf("compute" to 1)
        val root = googleCloudRootProject(testContextParameters())

        // With DefaultNumberOfBatches > 1, all packages get the feature except those
        // explicitly overridden to 1. We expect at least one package (compute) to be
        // missing the feature and the rest to keep it.
        val buildsWithoutFeature = collectPackageBuilds(root).filter { bt ->
            bt.features.items.filterIsInstance<ParallelTestsFeature>().isEmpty()
        }
        assertTrue("A batch override of 1 should disable the Parallel Tests feature for that service", buildsWithoutFeature.isNotEmpty())
    }
}
