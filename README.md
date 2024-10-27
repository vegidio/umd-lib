# Universal Media Downloader

**UMD** is a [Kotlin Multiplatform library](https://github.com/Kotlin/multiplatform-library-template) to easily extract links from media files hosted on popular websites.

It supports the following targets/platforms:

![](https://img.shields.io/badge/JVM-7F52FF?&style=for-the-badge&logo=kotlin&logoColor=white) ![](https://img.shields.io/badge/Android-34A853?style=for-the-badge&logo=android&logoColor=white) ![](https://img.shields.io/badge/iOS-FFFFFF?style=for-the-badge&logo=apple&logoColor=black) ![](https://img.shields.io/badge/macOS-000000?style=for-the-badge&logo=macos&logoColor=white) ![](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black) ![](https://img.shields.io/badge/Windows-0078D4?style=for-the-badge&logo=windows&logoColor=white) ![](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)

## ⬇️ Installation

This library can be installed in KMP / Android projects (through Maven), but also natively in other platforms such as iOS/macOS (Swift Package Manager) or Node.js (NPM).

### Maven (KMP / Android)

**UMD** is hosted in my own Maven repository, so before using it in your project you must add the repository `https://maven.vinicius.io` to your `settings.gradle.kts` file:

```kotlin
dependencyResolutionManagement {
    repositories {
        google()
        mavenCentral()
        maven("https://maven.vinicius.io")
    }
}
```

With the repository added, you just need to include the dependency in the file `build.gradle.kts`:

```kotlin
dependencies {
    implementation("io.vinicius.umd:umd:24.10.0")
}
```

### SwiftPM (iOS / macOS)

To add **UMD** to your Xcode project, select `File > Add Package Dependencies`:

![Xcode](docs/images/spm1.avif)

Enter the repository URL `https://github.com/vegidio/umd-lib` in the upper right corner to the screen and click on the button `Add Package`:

![Xcode](docs/images/spm2.avif)

### NPM (Node.js)

Coming soon...

## 🤖 Usage

Please visit the library's [website](https://vegidio.github.io/umd-lib) to find detailed instructions on how to use it in your project.

## 📝 License

**umd-lib** is released under the MIT License. See [LICENSE](LICENSE) for details.

## 👨🏾‍💻 Author

Vinicius Egidio ([vinicius.io](http://vinicius.io))
