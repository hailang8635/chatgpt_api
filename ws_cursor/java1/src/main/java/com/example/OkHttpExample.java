package com.example;

import okhttp3.*;
import java.io.IOException;

public class OkHttpExample {
    public static void main(String[] args) {
        // Create OkHttpClient instance
        OkHttpClient client = new OkHttpClient();

        // GET request example
        Request request = new Request.Builder()
            .url("https://news.baidu.com/")
            .build();

        // Synchronous request
        try {
            Response response = client.newCall(request).execute();
            if (response.isSuccessful()) {
                System.out.println("GET Request Successful!");
                System.out.println("Response Code: " + response.code());
                System.out.println("Content Length: " + response.body().contentLength());
                System.out.println("Content Type: " + response.body().contentType());
                System.out.println("\nFirst 200 characters of response:");
                String responseData = response.body().string();
                System.out.println(responseData.substring(0, Math.min(200, responseData.length())));
            }
        } catch (IOException e) {
            System.out.println("GET Request Failed: " + e.getMessage());
            e.printStackTrace();
        }

        // POST request example
        String json = "{\"name\":\"John\",\"age\":25}";
        RequestBody body = RequestBody.create(
            MediaType.parse("application/json; charset=utf-8"),
            json
        );

        Request postRequest = new Request.Builder()
            .url("https://api.example.com/post")
            .post(body)
            .build();

        try {
            Response response = client.newCall(postRequest).execute();
            System.out.println("\nPOST Request Result:");
            System.out.println("Success: " + response.isSuccessful());
            System.out.println("Response Code: " + response.code());
        } catch (IOException e) {
            System.out.println("\nPOST Request Failed: " + e.getMessage());
        }
    }
}
