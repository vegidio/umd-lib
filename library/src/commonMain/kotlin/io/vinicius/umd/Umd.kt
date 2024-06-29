package io.vinicius.umd

import co.touchlab.kermit.Logger
import co.touchlab.skie.configuration.annotations.DefaultArgumentInterop
import io.vinicius.umd.extractor.Extractor
import io.vinicius.umd.extractor.coomer.Coomer
import io.vinicius.umd.extractor.kemono.Kemono
import io.vinicius.umd.extractor.reddit.Reddit
import io.vinicius.umd.extractor.redgifs.Redgifs
import io.vinicius.umd.model.Event
import io.vinicius.umd.model.EventCallback
import io.vinicius.umd.model.ExtractorType
import io.vinicius.umd.model.Response
import io.vinicius.umd.util.Fetch

/**
 * A typealias for a map that represents metadata for different types of extractors.
 * The key is an [ExtractorType] which represents the type of the extractor.
 * The value is a map where the key is a string representing the metadata key and the value is any object representing
 * the metadata value.
 */
typealias Metadata = MutableMap<ExtractorType, Map<String, Any>>

/**
 * Umd class is responsible for handling media extraction from a given URL.
 *
 * @property url The URL from which media will be extracted.
 * @property metadata A mutable map containing metadata for different types of extractors.
 * @property callback An optional callback function that can be invoked during the extraction process.
 */
class Umd(
    private val url: String,
    private val metadata: Metadata = mutableMapOf(),
    val callback: EventCallback? = null,
) {
    private val tag = "Umd-Lib"
    private val extractor = findExtractor(url)

    /**
     * Queries the media found in the URL.
     *
     * @param limit The maximum number of files to query. Defaults to Int.MAX_VALUE.
     * @param extensions A list of file extensions to be queried. Defaults to an empty list.
     * @return A [Response] object containing information about the queried media.
     */
    @DefaultArgumentInterop.Enabled
    suspend fun queryMedia(limit: Int = Int.MAX_VALUE, extensions: List<String> = emptyList()): Response {
        val lowercaseExt = extensions.map { it.lowercase() }
        return extractor.queryMedia(url, limit, lowercaseExt)
    }

    /**
     * Configures a Fetch object to download media from the URL.
     *
     * @return A [Fetch] object configured for downloading media.
     */
    fun configureFetch(): Fetch = extractor.configureFetch()

    // region - Private methods
    private fun findExtractor(url: String): Extractor {
        val extractor = when {
            Coomer.isMatch(url) -> Coomer(callback = callback)
            Kemono.isMatch(url) -> Kemono(callback = callback)
            Reddit.isMatch(url) -> Reddit(callback = callback)

            Redgifs.isMatch(url) -> Redgifs(
                metadata = metadata[ExtractorType.RedGifs].orEmpty(),
                callback = callback,
            )

            else -> throw IllegalArgumentException("No extractor found for URL: $url")
        }

        val extractorName = extractor::class.simpleName?.lowercase().orEmpty()
        Logger.d(tag) { "Extractor found: $extractorName" }
        callback?.invoke(Event.OnExtractorFound(extractorName))

        return extractor
    }
    // endregion
}