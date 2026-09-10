/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package builds

import jetbrains.buildServer.configs.kotlin.BuildFeatures
import jetbrains.buildServer.configs.kotlin.buildFeatures.GolangFeature
import jetbrains.buildServer.configs.kotlin.buildFeatures.parallelTests

// NOTE: this file includes Extensions of the Kotlin DSL class BuildFeature
// This allows us to reuse code in the config easily, while ensuring the same build features can be used across builds.
// See the class's documentation: https://teamcity.jetbrains.com/app/dsl-documentation/root/build-feature/index.html


const val UseTeamCityGoTest = false

fun BuildFeatures.golang() {
    if (UseTeamCityGoTest) {
        feature(GolangFeature {
            testFormat = "json"
        })
    }
}

// parallelTests enables TeamCity's Parallel Tests build feature, which splits a
// build's tests into several batches that run in parallel on suitable agents and
// gathers the results back into a composite build overview.
// The feature is only enabled when numberOfBatches is greater than 1, so callers
// can safely pass the default (1) to keep the historical single-batch behaviour.
// See https://www.jetbrains.com/help/teamcity/parallel-tests.html
fun BuildFeatures.parallelTestsFeature(numberOfBatches: Int) {
    if (numberOfBatches > 1) {
        parallelTests {
            this.numberOfBatches = numberOfBatches
        }
    }
}
