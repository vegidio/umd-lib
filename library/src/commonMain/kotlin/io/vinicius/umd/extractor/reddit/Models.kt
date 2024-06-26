package io.vinicius.umd.extractor.reddit

import com.eygraber.uri.Url
import io.vinicius.umd.ktx.extension
import io.vinicius.umd.serializer.LocalDateTimeSerializer
import io.vinicius.umd.serializer.UrlSerializer
import kotlinx.datetime.Clock
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.TimeZone
import kotlinx.datetime.toLocalDateTime
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonObject

@Serializable
internal data class Submission(
    val data: Data,
) {
    @Serializable
    data class Data(
        val after: String?,
        val children: List<Child>,
    )
}

@Serializable
internal data class Child(
    val data: Data,
) {
    @Serializable
    data class Data(
        val author: String = "",

        @Serializable(UrlSerializer::class)
        val url: Url? = null,

        @Serializable(LocalDateTimeSerializer::class)
        val created: LocalDateTime = Clock.System.now().toLocalDateTime(TimeZone.UTC),

        @SerialName("is_gallery")
        val isGallery: Boolean = false,

        @SerialName("media_metadata")
        val mediaMetadata: JsonObject? = null,

        @SerialName("secure_media")
        val secureMedia: SecureMedia? = null,
    ) {
        val extension = url?.extension

        override fun equals(other: Any?): Boolean {
            val data = other as? Data
            return url == data?.url
        }

        override fun hashCode(): Int {
            var result = author.hashCode()
            result = 31 * result + url.hashCode()
            return result
        }
    }
}

@Serializable
internal data class MediaMetadata(
    val status: String,
    val s: S
) {
    @Serializable
    data class S(
        @SerialName("u")
        val image: String = "",
        val gif: String = ""
    )
}

@Serializable
internal data class SecureMedia(
    @SerialName("reddit_video")
    val redditVideo: RedditVideo? = null
) {
    @Serializable
    data class RedditVideo(
        @SerialName("fallback_url")
        val fallbackUrl: String,
    )
}