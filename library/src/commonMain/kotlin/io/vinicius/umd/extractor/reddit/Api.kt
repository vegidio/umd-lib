package io.vinicius.umd.extractor.reddit

import de.jensklingenberg.ktorfit.http.GET
import de.jensklingenberg.ktorfit.http.Path
import de.jensklingenberg.ktorfit.http.Query
import io.vinicius.umd.util.Fetch.Companion.ktorJson

internal interface Contract {
    @GET("{id}/.json?raw_json=1")
    suspend fun getSubmission(@Path("id") id: String): List<Submission>

    @GET("user/{user}/submitted.json?sort=new&raw_json=1")
    suspend fun getUserSubmissions(
        @Path("user") user: String,
        @Query("after") after: String,
        @Query("limit") limit: Int,
    ): Submission

    @GET("r/{subreddit}/hot.json?raw_json=1")
    suspend fun getSubredditSubmissions(
        @Path("subreddit") subreddit: String,
        @Query("after") after: String,
        @Query("limit") limit: Int,
    ): Submission
}

internal class RedditApi : Contract {
    private val api = ktorJson
        .baseUrl("https://www.reddit.com/")
        .build()
        .create<Contract>()

    override suspend fun getSubmission(id: String) = api.getSubmission(id)

    override suspend fun getUserSubmissions(user: String, after: String, limit: Int) =
        api.getUserSubmissions(user, after, limit)

    override suspend fun getSubredditSubmissions(subreddit: String, after: String, limit: Int) =
        api.getSubredditSubmissions(subreddit, after, limit)
}