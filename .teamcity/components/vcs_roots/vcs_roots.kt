/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package vcs_roots

import ProviderNameBeta
import ProviderNameGa
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot

object HashiCorpVCSRootGa: GitVcsRoot({
    name = "https://github.com/BBBmau/terraform-provider-${ProviderNameGa}#refs/heads/main"
    url = "https://github.com/BBBmau/terraform-provider-${ProviderNameGa}"
    branch = "refs/heads/main"
    branchSpec = """
        +:*
        -:refs/pull/*/head
    """.trimIndent()
})

object HashiCorpVCSRootBeta: GitVcsRoot({
    name = "https://github.com/BBBmau/terraform-provider-${ProviderNameBeta}#refs/heads/main"
    url = "https://github.com/BBBmau/terraform-provider-${ProviderNameBeta}"
    branch = "refs/heads/main"
    branchSpec = """
        +:*
        -:refs/pull/*/head
    """.trimIndent()
})

object ModularMagicianVCSRootGa: GitVcsRoot({
    name = "https://github.com/BBBmau/terraform-provider-${ProviderNameGa}#refs/heads/main"
    url = "https://github.com/BBBmau/terraform-provider-${ProviderNameGa}"
    branch = "refs/heads/main"
    branchSpec = """
        +:*
        -:refs/pull/*/head
    """.trimIndent()
})

object ModularMagicianVCSRootBeta: GitVcsRoot({
    name = "https://github.com/BBBmau/terraform-provider-${ProviderNameBeta}#refs/heads/main"
    url = "https://github.com/BBBmau/terraform-provider-${ProviderNameBeta}"
    branch = "refs/heads/main"
    branchSpec = """
        +:*
        -:refs/pull/*/head
    """.trimIndent()
})
