package io.vinicius.umd.extractor.kemono

internal sealed class SourceType {
    data class User(val service: String, val user: String) : SourceType()
    data class Post(val service: String, val user: String, val id: String) : SourceType()
}