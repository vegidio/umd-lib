package io.vinicius.umd.util

import co.touchlab.kermit.Logger
import co.touchlab.skie.configuration.annotations.DefaultArgumentInterop
import de.jensklingenberg.ktorfit.Ktorfit
import io.ktor.client.plugins.HttpRequestRetry
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.plugins.defaultRequest
import io.ktor.client.request.get
import io.ktor.client.request.header
import io.ktor.client.request.prepareGet
import io.ktor.client.statement.bodyAsChannel
import io.ktor.client.statement.bodyAsText
import io.ktor.client.utils.DEFAULT_HTTP_BUFFER_SIZE
import io.ktor.http.contentLength
import io.ktor.http.isSuccess
import io.ktor.serialization.kotlinx.json.json
import io.ktor.utils.io.core.isNotEmpty
import io.ktor.utils.io.core.readBytes
import io.vinicius.umd.exception.FetchException
import kotlinx.serialization.json.Json
import okio.Path.Companion.toPath
import okio.buffer
import okio.use

sealed class DownloadStatus {
    data class OnProgress(val bytes: Long, val total: Long) : DownloadStatus()
    data class OnComplete(val total: Long) : DownloadStatus()
}

private typealias DownloadCallback = (DownloadStatus) -> Unit

/**
 * Fetch class is used to send HTTP requests and download files.
 *
 * @property headers A map of headers to be included in the HTTP request.
 * @property retries The number of times to retry the request if it fails.
 * @property followRedirects A boolean indicating whether to follow redirects or not.
 */
class Fetch(
    headers: Map<String, String> = emptyMap(),
    retries: Int = 0,
    followRedirects: Boolean = true,
) {
    private val tag = "Fetch"

    private val ktorfit = Ktorfit.Builder()
        .httpClient {
            this.followRedirects = followRedirects
            defaultRequest {
                headers.entries.forEach { (key, value) ->
                    header(key, value)
                }
            }
            install(HttpRequestRetry) {
                maxRetries = retries
                retryIf { req, res ->
                    val shouldRetry = res.status.value == 429
                    if (shouldRetry) Logger.d(tag) { "#$retryCount Retrying request: ${req.url}" }
                    shouldRetry
                }
                exponentialDelay()
            }
        }
        .build()

    /**
     * Send a GET request to a URL and returns the response as string.
     *
     * @param url the URL that will receive the request.
     * @return the response represented as `String`.
     * @throws FetchException
     */
    suspend fun getString(url: String): String {
        val response = ktorfit.httpClient.get(url)

        return if (response.status.isSuccess()) {
            response.bodyAsText()
        } else {
            throw FetchException(response.status.value, response.bodyAsText())
        }
    }

    /**
     * Downloads a file.
     *
     * @param url the URL where the file is located.
     * @param filePath the path where the file will be saved.
     * @param callback lambda that tracks the download progress.
     */
    @DefaultArgumentInterop.Enabled
    suspend fun downloadFile(url: String, filePath: String, callback: DownloadCallback? = null) {
        ktorfit.httpClient.prepareGet(url).execute { response ->
            if (!response.status.isSuccess()) {
                throw FetchException(response.status.value, response.bodyAsText())
            }

            val fileSize = response.contentLength() ?: 0
            val channel = response.bodyAsChannel()
            val path = filePath.toPath()

            fileSystem.sink(path).buffer().use {
                while (!channel.isClosedForRead) {
                    val packet = channel.readRemaining(DEFAULT_HTTP_BUFFER_SIZE.toLong())
                    callback?.invoke(DownloadStatus.OnProgress(channel.totalBytesRead, fileSize))

                    while (packet.isNotEmpty) {
                        val bytes = packet.readBytes()
                        it.write(bytes)
                    }
                }
            }

            callback?.invoke(DownloadStatus.OnComplete(fileSize))
        }
    }

    companion object {
        internal val ktorJson = Ktorfit.Builder()
            .httpClient {
                install(ContentNegotiation) {
                    json(
                        Json {
                            prettyPrint = true
                            ignoreUnknownKeys = true
                        },
                    )
                }
            }
    }
}