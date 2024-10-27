package io.vinicius.umd.extractor.redgifs

import co.touchlab.kermit.Logger
import com.eygraber.uri.Url
import io.vinicius.umd.exception.InvalidSourceException
import io.vinicius.umd.extractor.Extractor
import io.vinicius.umd.model.Event
import io.vinicius.umd.model.EventCallback
import io.vinicius.umd.model.ExtractorType
import io.vinicius.umd.model.Media
import io.vinicius.umd.model.Response
import io.vinicius.umd.util.Fetch
import io.vinicius.umd.extractor.redgifs.Api as RedgifsApi

internal class Redgifs(
    private val api: RedgifsApi = RedgifsApi(),
    private val metadata: Map<String, Any> = emptyMap(),
    private val callback: EventCallback? = null,
) : Extractor {
    private val tag = this::class.simpleName.orEmpty()
    private val responseMeta = mutableMapOf<String, Any>()

    override suspend fun queryMedia(url: String, limit: Int, extensions: List<String>): Response {
        val source = getSourceType(url)

        val video = when (source) {
            is SourceType.Video -> fetchWatch(source)
        }

        val media = watchToMedia(video)
        callback?.invoke(Event.OnQueryCompleted(media.size))
        Logger.d(tag) { "Query completed: ${media.size} media found" }

        return Response(url, media, ExtractorType.RedGifs, responseMeta)
    }

    override fun configureFetch(): Fetch = Fetch(
        mapOf("User-Agent" to "UMD"),
    )

    // region - Private methods
    @Suppress("UseIfInsteadOfWhen")
    private fun getSourceType(url: String): SourceType {
        val regexVideo = """/watch/([^/\n?]+)""".toRegex()
        val id: String

        val source = when {
            url.contains(regexVideo) -> {
                val groups = regexVideo.find(url)?.groupValues
                id = groups?.get(1).orEmpty()
                SourceType.Video(id)
            }

            else -> throw InvalidSourceException(url, ExtractorType.RedGifs, "No support for URL: $url")
        }

        val sourceName = source::class.simpleName?.lowercase().orEmpty()
        callback?.invoke(Event.OnExtractorTypeFound(sourceName, id))
        Logger.d(tag) { "Extractor type found: $sourceName" }

        return source
    }

    private suspend fun fetchWatch(source: SourceType.Video): Watch {
        val token = if (metadata.containsKey("token")) {
            Logger.d(tag) { "Reusing previous token" }
            metadata["token"] as String
        } else {
            Logger.d(tag) { "Requesting new token" }
            api.getToken().token
        }

        responseMeta["token"] = token

        val video = api.getVideo(
            token = "Bearer $token",
            videoUrl = "https://www.redgifs.com/watch/${source.id}",
            videoId = source.id,
        )

        callback?.invoke(Event.OnMediaQueried(1))
        Logger.d(tag) { "Media queried: 1" }

        return video
    }

    private fun watchToMedia(watch: Watch): List<Media> {
        return listOf(
            Media(
                watch.gif.url.hd,
                mapOf(
                    "source" to "watch",
                    "name" to watch.gif.userName,
                    "created" to watch.gif.created.toString(),
                    "duration" to watch.gif.duration,
                    "id" to watch.gif.id,
                ),
            ),
        )
    }
    // endregion

    companion object {
        fun isMatch(url: String): Boolean {
            val urlObj = Url.parse(url)
            return urlObj.host.endsWith("redgifs.com", true)
        }
    }
}