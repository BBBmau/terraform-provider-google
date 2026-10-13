/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

import generated.PackagesListBeta
import generated.PackagesListGa
import generated.ServicesListBeta
import generated.ServicesListGa
import projects.getPackageNameList

// Each service package gets its own exclusive (quota 1) shared resource per test project instead of a value in a
// custom-value resource. TeamCity holds quota locks on a Parallel Tests composite for all of its batches, but it
// cannot do that for custom-value locks, so this keeps one lock per service across every batch.
fun serviceLockName(sharedResourceName: String, packageName: String): String = "$sharedResourceName - $packageName"

fun serviceLockPackageNames(sharedResourceName: String): List<String> = when (sharedResourceName) {
    SharedResourceNameGa -> getPackageNameList(ServicesListGa) + getPackageNameList(PackagesListGa)
    // Use Beta list of services here, assuming Beta is a superset of GA
    SharedResourceNameBeta, SharedResourceNameVcr -> getPackageNameList(ServicesListBeta) + getPackageNameList(PackagesListBeta)
    else -> throw IllegalArgumentException("unknown shared resource $sharedResourceName")
}

fun serviceLockNames(sharedResourceName: String): List<String> =
    serviceLockPackageNames(sharedResourceName).map { serviceLockName(sharedResourceName, it) }
