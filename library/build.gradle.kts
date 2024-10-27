import org.jetbrains.kotlin.gradle.ExperimentalKotlinGradlePluginApi
import org.jetbrains.kotlin.gradle.dsl.JvmTarget
import org.jetbrains.kotlin.gradle.plugin.mpp.apple.XCFramework

plugins {
    alias(libs.plugins.android.library)
    alias(libs.plugins.kotlin.multiplatform)
    alias(libs.plugins.kotlin.serialization)
    alias(libs.plugins.skie)
    alias(libs.plugins.ktlint)

    `maven-publish`
}

kotlin {
    applyDefaultHierarchyTemplate()

    // JVM
    jvm()

    // Android
    androidTarget {
        publishLibraryVariants("release")
        mavenPublication { artifactId = "android" }

        @OptIn(ExperimentalKotlinGradlePluginApi::class)
        compilerOptions {
            jvmTarget.set(JvmTarget.JVM_21)
        }
    }

    // iOS & macOS
    val frameworkName = "UMD"
    val xcf = XCFramework(frameworkName)
    val appleTargets = listOf(iosArm64(), iosSimulatorArm64(), iosX64(), macosArm64(), macosX64())

    appleTargets.forEach {
        it.binaries.framework {
            baseName = frameworkName
            binaryOption("bundleId", "io.vinicius.umd")
            xcf.add(this)
        }
    }

    // Linux
    linuxArm64()
    linuxX64()

    // Windows
    mingwX64()

    sourceSets {
        // Common
        commonMain.dependencies {
            implementation(libs.coroutines.core)
            implementation(libs.kermit)
            implementation(libs.klopik)
            implementation(libs.kotlin.datetime)
            implementation(libs.kotlin.serialization)
            implementation(libs.ksoup)
            implementation(libs.okio)
            implementation(libs.skie.annotations)
            implementation(libs.slf4j)
            implementation(libs.uri)
        }

        commonTest.dependencies {
            implementation(libs.coroutines.test)
            implementation(libs.kotlin.test)
        }

        all {
            languageSettings.optIn("kotlin.js.ExperimentalJsExport")
        }
    }
}

android {
    namespace = "io.vinicius.umd"
    compileSdk = libs.versions.android.compileSdk.get().toInt()
    defaultConfig {
        minSdk = libs.versions.android.minSdk.get().toInt()
    }
}

// Disable Skie analytics
skie {
    analytics {
        enabled.set(false)
    }
}

// Workaround for KSP picking the wrong Java version
afterEvaluate {
    tasks.withType<JavaCompile>().configureEach {
        sourceCompatibility = JavaVersion.VERSION_21.toString()
        targetCompatibility = JavaVersion.VERSION_21.toString()
    }
}

configure<org.jlleitschuh.gradle.ktlint.KtlintExtension> {
    additionalEditorconfig.set(
        mapOf("ktlint_code_style" to "intellij_idea"),
    )
    filter {
        exclude {
            it.file.path.contains("generated")
        }
    }
}

group = "io.vinicius.umd"
version = System.getenv("VERSION") ?: "1.0-SNAPSHOT"

afterEvaluate {
    apply(from = "../publish.gradle.kts")
}