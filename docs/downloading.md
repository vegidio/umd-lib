# Downloading Files

You can download most media files using any browser or tools like `wget` and `curl`. However, if you intend to use **UMD** in your project, you can use the class `Fetch` if you don't want to write your own code to download the files:

=== "Kotlin"

    ```kotlin linenums="1" hl_lines="4 6"
    val umd = Umd("https://www.reddit.com/user/atomicbrunette18")
    val response = umd.queryMedia()

    val fetch = umd.configureFetch()
    response.media.forEachIndexed { index, media ->
        fetch.downloadFile(media.url, "file${index}.${media.extension}")
    }
    ```

=== "Swift"

    ```swift linenums="1" hl_lines="4 6"
    let umd = Umd(url: "https://www.reddit.com/user/atomicbrunette18")
    let response = try! await umd.queryMedia()

    let fetch = umd.queryMedia()
    for (index, media) in response.media.enumerated() {
        try? await fetch.downloadFile(url: media.url, filePath: "file\(index).\(media.extension!)")
    }
    ```

=== "TypeScript"

    Coming soon...

!!! warning "Attention"

    The code above is just a quick and dirty example on how the use the library but it's not production ready. In other words, when you use this library in your project make sure to always catch exceptions, check for nullability and follow other best practices.

## Special Cases

While most files can be downloaded without restrictions, some websites create mechanisms to prevent the media to be downloaded, even when you have the direct link to the file.

If you use the `Fetch` object, like described in the example above, you should be fine since `Fetch` will automatically do whatever is necessary to circumvent these restrictions, but if you want to download the files by your own means then follow the instructions below:

### Coomer

If you're downloading a small number of files then you should not have any problem, however at some point Coomer will stop accepting requests to download the media and it will return an error `HTTP 429 Too Many Request`.

To circumvent this, when you receive an `HTTP 429` error, it's recommended to wait **at least 15 seconds** before trying to download anything again.

### Kemono

If you're downloading a small number of files then you should not have any problem, however at some point Kemono will stop accepting requests to download the media and it will return an error `HTTP 429 Too Many Request`.

To circumvent this, when you receive an `HTTP 429` error, it's recommended to wait **at least 15 seconds** before trying to download anything again.

### RedGifs

RedGifs protection requires that the client downloading the file has the same user agent of the client that queried the media URL. That means that if you query the media URL with **UMD** and tries to open it in a browser, for example, it won't work since **UMD** and the browser have different user agents.

To circumvent this, when you try to download a RedGifs media you must make sure to set the user agent HTTP header to `UMD`, like this: `User-Agent: UMD`.
