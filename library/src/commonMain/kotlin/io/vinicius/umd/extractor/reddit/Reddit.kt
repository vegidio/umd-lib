package io.vinicius.umd.extractor.reddit

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
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.decodeFromJsonElement
import io.vinicius.umd.extractor.reddit.Api as RedditApi

internal class Reddit(
    private val api: RedditApi = RedditApi(),
    private val callback: EventCallback? = null,
) : Extractor {
    private val tag = this::class.simpleName.orEmpty()

    override suspend fun queryMedia(url: String, limit: Int, extensions: List<String>): Response {
        var sourceName = ""
        val source = getSourceType(url)

        val submissions = when (source) {
            is SourceType.Submission -> {
                sourceName = source.name
                fetchSubmissions(source, limit, extensions)
            }

            is SourceType.User -> {
                sourceName = source.name
                fetchSubmissions(source, limit, extensions)
            }

            is SourceType.Subreddit -> {
                sourceName = source.name
                fetchSubmissions(source, limit, extensions)
            }
        }

        val media = submissionsToMedia(submissions, source, sourceName)
        callback?.invoke(Event.OnQueryCompleted(media.size))
        Logger.d(tag) { "Query completed: ${media.size} media found" }

        return Response(url, media, ExtractorType.Reddit, emptyMap())
    }

    override fun configureFetch(): Fetch = Fetch()

    // region - Private methods
    private fun getSourceType(url: String): SourceType {
        val regexSubmission = """/(?:r|u|user)/([^/?]+)/comments/([^/\n?]+)""".toRegex()
        val regexUser = """/(?:u|user)/([^/\n?]+)""".toRegex()
        val regexSubreddit = """/r/([^/\n]+)""".toRegex()
        val name: String

        val source = when {
            url.contains(regexSubmission) -> {
                val groups = regexSubmission.find(url)?.groupValues
                name = groups?.get(1).orEmpty()
                val id = groups?.get(2).orEmpty()
                SourceType.Submission(name, id)
            }

            url.contains(regexUser) -> {
                val groups = regexUser.find(url)?.groupValues
                name = groups?.get(1).orEmpty()
                SourceType.User(name)
            }

            url.contains(regexSubreddit) -> {
                val groups = regexSubreddit.find(url)?.groupValues
                name = groups?.get(1).orEmpty()
                SourceType.Subreddit(name)
            }

            else -> throw InvalidSourceException(url, ExtractorType.Reddit, "No support for URL: $url")
        }

        val sourceName = source::class.simpleName?.lowercase().orEmpty()
        callback?.invoke(Event.OnExtractorTypeFound(sourceName, name))
        Logger.d(tag) { "Extractor type found: $sourceName" }

        return source
    }

    private suspend fun fetchSubmissions(
        source: SourceType,
        limit: Int,
        extensions: List<String>,
    ): List<Child> {
        val submissions = mutableSetOf<Child>()
        var after: String? = ""

        do {
            val response = when (source) {
                is SourceType.Submission -> api.getSubmission(source.id).first()
                is SourceType.User -> api.getUserSubmissions(source.name, after.orEmpty(), 100)
                is SourceType.Subreddit -> api.getSubredditSubmissions(source.name, after.orEmpty(), 100)
            }

            val filteredSubmissions = response.data.children
                .flatMap { getGallerySubmissions(it) }
                .filter { extensions.isEmpty() || extensions.contains(it.data.extension) }

            after = response.data.after
            val amountBefore = submissions.size
            submissions.addAll(filteredSubmissions)

            val queried = submissions.size - amountBefore
            callback?.invoke(Event.OnMediaQueried(queried))
            Logger.d(tag) { "Media queried: $queried" }
        } while (response.data.children.isNotEmpty() && submissions.size < limit && after != null)

        return submissions.take(limit)
    }

    private fun getGallerySubmissions(child: Child): List<Child> {
        val json = Json { ignoreUnknownKeys = true }

        return if (child.data.isGallery) {
            val jsonObject = child.data.mediaMetadata

            jsonObject?.keys?.mapNotNull {
                val mm = json.decodeFromJsonElement<MediaMetadata>(jsonObject.getValue(it))

                if (mm.status == "valid") {
                    val url = mm.s.image.ifBlank { mm.s.gif }
                    child.copy(data = Child.Data(url = Url.parse(url), isGallery = true))
                } else {
                    null
                }
            }.orEmpty()
        } else {
            listOf(child)
        }
    }

    private fun submissionsToMedia(submissions: List<Child>, source: SourceType, name: String): List<Media> {
        return submissions.map {
            val url = it.data.secureMedia?.redditVideo?.fallbackUrl ?: it.data.url.toString()

            Media(
                url,
                mapOf(
                    "source" to source::class.simpleName?.lowercase(),
                    "name" to name,
                    "created" to it.data.created.toString(),
                ),
            )
        }
    }
    // endregion

    companion object {
        fun isMatch(url: String): Boolean {
            val urlObj = Url.parse(url)
            return urlObj.host.endsWith("reddit.com", true)
        }
    }
}