package io.vinicius.umd

import co.touchlab.kermit.Logger
import io.vinicius.umd.extractor.Extractor
import io.vinicius.umd.extractor.coomer.Coomer
import io.vinicius.umd.extractor.reddit.Reddit
import io.vinicius.umd.extractor.redgifs.Redgifs
import io.vinicius.umd.model.Event
import io.vinicius.umd.model.EventCallback
import io.vinicius.umd.model.ExtractorType
import io.vinicius.umd.util.Fetch
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.promise

typealias Metadata = MutableMap<ExtractorType, Map<String, Any>>

internal val scope = CoroutineScope(Dispatchers.Default)

@JsExport
class Umd(
    private val url: String,
    private val metadata: Map<ExtractorType, Map<String, Any>> = mutableMapOf(),
    val callback: EventCallback? = null,
) {
    private val tag = "Umd-Lib"
    private val extractor = findExtractor(url)

    /**
     * Query the media found in the URL.
     *
     * @param limit the max amount of files that you want to query.
     * @param extensions list of file extensions that you want to be queried.
     * @return `Response` object with information about the media queried.
     */
    fun queryMedia(limit: Int = Int.MAX_VALUE, extensions: Array<String> = emptyArray()) = scope.promise {
        val lowercaseExt = extensions.map { it.lowercase() }
        extractor.queryMedia(url, limit, lowercaseExt)
    }

    /**
     * Gets the properly configured Fetch object to download media from this URL.
     *
     * @return `Fetch` object.
     */
    fun configureFetch(): Fetch = extractor.configureFetch()

    // region - Private methods
    private fun findExtractor(url: String): Extractor {
        val extractor = when {
            Coomer.isMatch(url) -> Coomer(callback = callback)
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