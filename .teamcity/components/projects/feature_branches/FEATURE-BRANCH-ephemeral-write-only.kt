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

    val trigger  = NightlyTriggerConfiguration(
        branch = "refs/heads/$featureBranchEphemeralWriteOnly" // triggered builds must test the feature branch
    )


    // GA
    val gaConfig = getGaAcceptanceTestConfig(allConfig)
    // These are the packages that have resources that will use write-only attributes
    var ServicesListWriteOnlyGa = mapOf(
        "compute" to mapOf(
            "name" to "compute",
            "displayName" to "Compute",
            "path" to "./google/services/compute"
        ),
        "secretmanager" to mapOf(
            "name" to "secretmanager",
            "displayName" to "Secretmanager",
            "path" to "./google/services/secretmanager"
        ),
        "sql" to mapOf(
            "name" to "sql",
            "displayName" to "Sql",
            "path" to "./google/services/sql"
        ),
        "bigquery_datatransfer" to mapOf(
            "name" to "bigquery_datatransfer",
            "displayName" to "Bigquery Datatransfer",
            "path" to "./google/services/bigquery_datatransfer"
        )
    )

    val projectId = replaceCharsId(featureBranchEphemeralWriteOnly)
    val buildConfigsGa = BuildConfigurationsForPackages(ServicesListWriteOnlyGa, ProviderNameGa, projectId, HashiCorpVCSRootGa, listOf(SharedResourceNameGa), gaConfig)

    // Beta
    val betaConfig = getBetaAcceptanceTestConfig(allConfig)
    var ServicesListWriteOnlyBeta = mapOf(
        "compute" to mapOf(
            "name" to "compute",
            "displayName" to "Compute - Beta",
            "path" to "./google-beta/services/compute"
        ),
        "secretmanager" to mapOf(
            "name" to "secretmanager",
            "displayName" to "Secretmanager - Beta",
            "path" to "./google-beta/services/secretmanager"
        ),
        "sql" to mapOf(
            "name" to "sql",
            "displayName" to "Sql - Beta",
            "path" to "./google-beta/services/sql"
        ),
        "bigquery_datatransfer" to mapOf(
            "name" to "bigquery_datatransfer",
            "displayName" to "Bigquery Datatransfer - Beta",
            "path" to "./google-beta/services/bigquery_datatransfer"
        )
    )
    val projectIdBeta = replaceCharsId(featureBranchEphemeralWriteOnly + "-beta")
    val buildConfigsBeta = BuildConfigurationsForPackages(ServicesListWriteOnlyBeta, ProviderNameBeta, projectIdBeta, HashiCorpVCSRootBeta, listOf(SharedResourceNameBeta), betaConfig)

    // Make all builds use a 1.10.0-ish version of TF core
    val allBuildConfigs = buildConfigsGa + buildConfigsBeta
    allBuildConfigs.forEach{ builds ->
        builds.overrideTerraformCoreVersion(EphemeralWriteOnlyTfCoreVersion)
    }

    // ------

    return Project{
        id(projectId)
        name = featureBranchEphemeralWriteOnly
        description = "Subproject for testing feature branch $featureBranchEphemeralWriteOnly"

        // Register all build configs in the project
        allBuildConfigs.forEach{ builds ->
            buildType(builds)
        }

        params {
            readOnlySettings()
        }
    }
}