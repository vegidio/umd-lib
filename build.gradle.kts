plugins {
    alias(libs.plugins.android.library) apply false
    alias(libs.plugins.kotlin.multiplatform) apply false
    alias(libs.plugins.detekt)
    alias(libs.plugins.ktlint)
}

detekt {
    config.setFrom("$rootDir/config/detekt.yml")
    source.setFrom(
        "$rootDir/library/src/commonMain/kotlin",
        "$rootDir/library/src/commonTest/kotlin",
        "$rootDir/library/src/jsMain/kotlin",
        "$rootDir/library/src/nativeMain/kotlin",
        "$rootDir/library/src/nativeTest/kotlin",
    )
}