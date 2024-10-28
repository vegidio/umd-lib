// swift-tools-version: 5.7
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "UMD",
    platforms: [
        .iOS(.v16),
        .macOS(.v13)
    ],
    products: [
        .library(name: "UMD", targets: ["UMD"])
    ],
    targets: [
        .binaryTarget(
            name: "UMD",
            url: "https://github.com/vegidio/umd-lib/releases/download/24.10.1/umd-xcframework.zip",
            checksum: "b904f2f40402b34f20cb5d947b80b0df9812f5ee9e9fdddc0c4f2cf11ae0eea8"
        )
    ]
)
