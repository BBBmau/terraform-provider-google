/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package generated

var ServiceParallelism = mapOf(
    "looker" to 1
)

// ServiceNumberOfBatches optionally overrides the number of TeamCity Parallel
// Tests batches used for a specific service package. Packages not listed here
// fall back to DefaultNumberOfBatches. A value of 1 keeps the package on a
// single agent (Parallel Tests disabled for that package).
var ServiceNumberOfBatches = mapOf<String, Int>()
