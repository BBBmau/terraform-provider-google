/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects

import GlobalSweepersProjectName
import DefaultBranchName
import NightlyTestsProjectId
import ServiceSweeperName
import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import builds.*
import generated.SweepersListGa
import jetbrains.buildServer.configs.kotlin.AbsoluteId
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.BuildTypeSettings
import jetbrains.buildServer.configs.kotlin.DslContext
import jetbrains.buildServer.configs.kotlin.FailureAction
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.triggers.finishBuildTrigger
import replaceCharsId
import vcs_roots.HashiCorpVCSRootGa

// globalSweepersSubProject returns a subproject that contains sweepers for global resources (projects, folders)
// Sweeping projects is an edge case because it doesn't respect boundaries between different testing projects GA/Beta/PR
fun globalSweepersSubProject(allConfig: AllContextParameters): Project {

    val sweeperId = replaceCharsId("GLOBAL_SWEEPER")

    // Get config for using the GA identity (arbitrary choice as sweeper isn't confined by GA/Beta etc.)
    val gaConfig = getGaAcceptanceTestConfig(allConfig)

    // List of ALL shared resources; avoid clashing with any other running build
    val sharedResources: List<String> = listOf(SharedResourceNameGa, SharedResourceNameBeta, SharedResourceNameVcr)

    // Match the GA service sweeper ID created by googleSubProjectGa() and nightlyTests().
    val gaProjectId = replaceCharsId("GOOGLE")
    val betaProjectId = replaceCharsId("GOOGLE_BETA")
    val gaNightlyTestsId = "${DslContext.projectId}_${replaceCharsId("${gaProjectId}_${NightlyTestsProjectId}_all_tests")}"
    val gaServiceSweeperId = "${DslContext.projectId}_${replaceCharsId("${gaProjectId}_${NightlyTestsProjectId}_${ServiceSweeperName}")}"
    val betaServiceSweeperId = "${DslContext.projectId}_${replaceCharsId("${betaProjectId}_${NightlyTestsProjectId}_${ServiceSweeperName}")}"

    // Join both provider sweeper chains before global cleanup. Manual runs of the
    // global sweepers do not depend on this gate and therefore do not start tests.
    val nightlySweeperGate = BuildType {
        id(replaceCharsId("${sweeperId}_NIGHTLY_SWEEPER_GATE"))
        name = "Nightly Sweeper Gate"
        type = BuildTypeSettings.Type.COMPOSITE
        triggers {
            finishBuildTrigger {
                buildType = gaNightlyTestsId
                branchFilter = "+:$DefaultBranchName"
                successfulOnly = false
            }
        }
        dependencies {
            snapshot(AbsoluteId(gaServiceSweeperId)) {
                onDependencyFailure = FailureAction.IGNORE
                onDependencyCancel = FailureAction.IGNORE
            }
            snapshot(AbsoluteId(betaServiceSweeperId)) {
                onDependencyFailure = FailureAction.IGNORE
                onDependencyCancel = FailureAction.IGNORE
            }
        }
    }

    // Create build config for sweeping project resources
    // Uses the HashiCorpVCSRootGa VCS Root so that the latest sweepers in hashicorp/terraform-provider-google are used
    val projectSweeperConfig = BuildConfigurationForGlobalSweeper("N/A", "Project Sweeper", "GoogleProject", SweepersListGa, sweeperId, HashiCorpVCSRootGa, sharedResources, gaConfig)
    // Create build config for sweeping folder resources
    val folderSweeperConfig = BuildConfigurationForGlobalSweeper("N/A", "Folder Sweeper", "GoogleFolder", SweepersListGa, sweeperId, HashiCorpVCSRootGa, sharedResources, gaConfig)
    val sweepers = listOf(projectSweeperConfig, folderSweeperConfig)
    sweepers.forEach { sweeper ->
        sweeper.triggers {
            finishBuildTrigger {
                buildType = nightlySweeperGate.id!!.value
                branchFilter = "+:$DefaultBranchName"
                successfulOnly = false
            }
        }
    }

    return Project{
        id(sweeperId)
        name = GlobalSweepersProjectName
        description = "Subproject containing build configurations for sweeping global resources like projects and folders"

        // Register build configs in the project
        buildType(nightlySweeperGate)
        sweepers.forEach { buildType(it) }

        params {
            readOnlySettings()
        }
    }
}