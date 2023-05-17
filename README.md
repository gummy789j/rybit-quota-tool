# Rybit Quota Calculation Tool

## Motivation
Because the Rybit deposit quota calculation is based on a **rolling cumulative** method, it becomes challenging to determine your remaining quota. Therefore, we have developed a tool that enables quick calculation of your remaining quota. check the official [document](https://www.rybit.com/blog/member-level) here.
    

## Example 

Assume that I am VIP level 2 on Rybit, according to the calculated result, I still can deposit 500,000 TWD in 24 hours and 3,700,000 TWD in 30 days.

The following table of unlock time & amount represents that your 30 days quota will be unlocked at 2023-05-19 17:08:50 UTC+8 for 30,000 TWD and so on.

```
--------------------------------------------------
*** current remaining deposit quota in 24h: 500000
--------------------------------------------------
*** current remaining deposit quota in 30 days: 3700000
------------------------  -----------
            unlock time       amount 
------------------------  -----------
    2023-05-19 17:08:50        30000 

    2023-06-09 12:35:48       150000 

    2023-06-10 12:50:38       150000 

    2023-06-11 13:20:08       150000 

    2023-06-14 12:37:18       320000 
------------------------  -----------
```

## Quick Start

- step 1: execute the git command below on your target folder 
    
    ```
    git clone https://github.com/gummy789j/rybit-quota-tool.git
    ```
    tips: please install [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) before you started

- step 2: setting up your AUTH_TOKEN and VIP_LEVEL in Makefile
    
    ```
    AUTH_TOKEN := Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJiaXR1bml2ZXJzZS5vcmcvYXBpLXRyYWRpbmcvdHJhZGUiLCJleHAiOjE2ODQzNDQ2NjUsImp0aSI6IjIzNTE0OGFiLWUwYmItNDQzOC04ZGFlLKJNDBQ3ODIzNjAyYSIsImlhdCI6MTY4NDM0Mjg2NSwiaXNzIjoiYml0dW5pdmVyc2Uub3JnL3RyYWRlL3NpZ25lciIsInN1YiI6ImJpdHVuaXZlcnNlLm9yZy9hcGktdHJhZGluZy90cmFkZSIsInZlcnNpb24iOiIxLjAiLCJkYXRhIjp7InVzZXJfaWQiOiIyODc5NzgwOC1kMmMyLTRkOGUtOGUyYS00YWNmOGJjMmM1NjMiLCJrZXlfaWQiOiIwODhmZGI4YWIyZDNkNDVjZjcxMzI3NjM0MjMwNGE2OCIsImV4Y2hhbmdlIjoicGlvbmV4LnYyIiwiYXBpX2tleSI6ImM0YWM0ZGMwLWUyNTAtNGNlZi04ZTU1LTFmYTg5MWZkYTkwMSIsImV4Y2hhbmdlX2luZm8iOiJib3QifX0.iE9oECmT1pY0f50i0wUM0uVIHuQVrIXgHl-y3nN9gt2Gepk3-TsV-bB-5JOID4hQiOwXlAqkOz_lTLJlZDSX8g

    VIP_LEVEL := 2
    ```

- step 3: according your OS, executing command below 
    - macOS
    ```
    make rqt-mac
    ```
    - linux
    ```
    make rqt-linux
    ```
    - windows
    ```
    make rqt-windows
    ```


