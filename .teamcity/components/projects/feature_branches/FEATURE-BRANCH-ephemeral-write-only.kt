/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects.feature_branches

import ProviderNameBeta
import ProviderNameGa
import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import builds.*
import generated.ServicesListBeta
import generated.ServicesListGa
import jetbrains.buildServer.configs.kotlin.Project
import replaceCharsId
import vcs_roots.HashiCorpVCSRootBeta
import vcs_roots.HashiCorpVCSRootGa
import vcs_roots.ModularMagicianVCSRootBeta
import vcs_roots.ModularMagicianVCSRootGa

const val featureBranchEphemeralWriteOnly = "FEATURE-BRANCH-ephemeral-write-only"
const val EphemeralWriteOnlyTfCoreVersion = "1.10.0"

fun featureBranchEphemeralWriteOnlySubProject(allConfig: AllContextParameters): Project {

    val projectId = replaceCharsId(featureBranchEphemeralWriteOnly)

    val packageName = "compute" 
    val vcrConfig = getVcrAcceptanceTestConfig(allConfig) // Reused below for both MM testing build configs
    val trigger  = NightlyTriggerConfiguration(
        branch = "refs/heads/$featureBranchEphemeralWriteOnly" // triggered builds must test the feature branch
    )


    // GA
    val gaConfig = getGaAcceptanceTestConfig(allConfig)
    // How to make only build configuration to the relevant package(s)
    val resourceManagerPackageGa = ServicesListGa.getValue(packageName)

    // Enable testing using hashicorp/terraform-provider-google
    var parentId = "${projectId}_HC_GA"
    val buildConfigHashiCorpGa = BuildConfigurationForSinglePackage(packageName, resourceManagerPackageGa.getValue("path"), "Ephemeral Write Only in $packageName (GA provider, HashiCorp downstream)", ProviderNameGa, parentId, HashiCorpVCSRootGa, listOf(SharedResourceNameGa), gaConfig)
    buildConfigHashiCorpGa.addTrigger(trigger)

    // Enable testing using modular-magician/terraform-provider-google
    parentId = "${projectId}_MM_GA"
    val buildConfigModularMagicianGa = BuildConfigurationForSinglePackage(packageName, resourceManagerPackageGa.getValue("path"), "Ephemeral Write Only in $packageName (GA provider, MM upstream)", ProviderNameGa, parentId, ModularMagicianVCSRootGa, listOf(SharedResourceNameVcr), vcrConfig)
    // No trigger added here (MM upstream is manual only)

    // Beta
    val betaConfig = getBetaAcceptanceTestConfig(allConfig)
    val resourceManagerPackageBeta = ServicesListBeta.getValue(packageName)

    // Enable testing using hashicorp/terraform-provider-google-beta
    parentId = "${projectId}_HC_BETA"
    val buildConfigHashiCorpBeta = BuildConfigurationForSinglePackage(packageName, resourceManagerPackageBeta.getValue("path"), "Ephemeral Write Only in $packageName (Beta provider, HashiCorp downstream)", ProviderNameBeta, parentId, HashiCorpVCSRootBeta, listOf(SharedResourceNameBeta), betaConfig)
    buildConfigHashiCorpBeta.addTrigger(trigger)

    // Enable testing using modular-magician/terraform-provider-google-beta
    parentId = "${projectId}_MM_BETA"
    val buildConfigModularMagicianBeta = BuildConfigurationForSinglePackage(packageName, resourceManagerPackageBeta.getValue("path"), "Ephemeral Write Only in $packageName (Beta provider, MM upstream)", ProviderNameBeta, parentId, ModularMagicianVCSRootBeta, listOf(SharedResourceNameVcr), vcrConfig)
    // No trigger added here (MM upstream is manual only)


    // ------

    // Make all builds use a 1.10.0-ish version of TF core
    val allBuildConfigs = listOf(buildConfigHashiCorpGa, buildConfigModularMagicianGa, buildConfigHashiCorpBeta, buildConfigModularMagicianBeta)
    allBuildConfigs.forEach{ b ->
        b.overrideTerraformCoreVersion(EphemeralWriteOnlyTfCoreVersion)
    }

    // ------

    return Project{
        id(projectId)
        name = featureBranchEphemeralWriteOnly
        description = "Subproject for testing feature branch $featureBranchEphemeralWriteOnly"

        // Register all build configs in the project
        allBuildConfigs.forEach{ b ->
            buildType(b)
        }

        params {
            readOnlySettings()
        }
    }
}