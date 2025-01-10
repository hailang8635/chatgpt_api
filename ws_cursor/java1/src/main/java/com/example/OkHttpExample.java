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
                String responseData = response.body().string();
                System.out.println("Response data: " + responseData);
            }
        } catch (IOException e) {
            e.printStackTrace();
        }

        // Asynchronous request
        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                e.printStackTrace();
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    String responseData = response.body().string();
                    System.out.println("Async response data: " + responseData);
                }
            }
        });

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
            if (response.isSuccessful()) {
                System.out.println("POST request successful");
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
