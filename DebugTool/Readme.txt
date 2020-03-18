DebugTool.pdf
1) REST
    1) Create HTML Server lisening on 8080
    2) Create webservices with callback handler
        / : HandleMainCOnfig : THis Handler generates HTML page back to
            client (brower). We can display Form in that HMTL page. On
            Form submit - This will generate one more req to rest server
            for differnt webserice say /debugsubmit
