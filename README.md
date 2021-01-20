# FQimg
## fqimg playground
http://loveczp.github.io/fqimg/  

[Chinese Version document](https://github.com/loveczp/fqimg/blob/master/README-cn.md)

FQimg is a image server powered by Golang  
features:

* dynamicly process the images。The url of the image can accept commands and parameters, by which we can crop, fill, grayscale,etc the image in real time. 
* chain mode。commands can be chained together like Unix pipe. 
* high performance. Because of the excelence of Golang, the FQimg server has very low RAM and CPU consumption.
* support webp webp image format. Webp  is better than jpeg, it has the lower size on the same image quality. For mobile application, it is very important.
* simple to deploy. FQimg is an executalbe binary file, don't need to install any extra library. 
* support multiple image format, such as jpeg, gif,bmp, webp. Also support the quality parameter for each format.
* FQimg has cache mechanism, repeatly requested images are cached. This feature can tremendously increase the performance. 
* support the uploading Ip control. Only designated IPs are allowed to upload image.
* CORS upload support. This is important for web application.
* support HTTPS
* support two storage type. 1. local file system，2. [seaweeddfs](https://github.com/chrislusf/seaweedfs)


## examples

* **original image**
 https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0

* **dynamic cut example**  
 fill=400_400 command will dynamicly cut the image in to 400*400 image, the height and width can be dynamicly set at one's will.
 https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=400_400

* **china mode example**  
 fit=200_300&grayscale is two chained command, fit and grascale. the server would first cut the image by 200*300, and then grayscale the product of the first step. Surely we can chain more command.  
 https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fit=200_300&grayscale

* **output webp format**  
 50 is the quality parameter of webp command  
 https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?webp=50





## Document


* **install and run**  
go get github.com/loveczp/fqimg  
then run fqimg with following command  
  ```bash
  fqimg -c=path/to/config/file
  ```
    attention:
    * on window, we should first install [tdm-gcc](http://tdm-gcc.tdragon.net/download), or FQimg    * can not be compiled. Because the webp library is a C library.
    * go version should be greater than 1.8
* **image upload**  
 we can use the curl command to post a image to the server  
  ```bash
  curl -F "file=@xxxxx.jpg" "http://fqimg.com/put"  
   ```
   
    change the xxx.jpg and fqimg.com according to your circumstance.  
    ["https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0"]  
    when local storage is used, the last part of the URL is the of MD5 value of the image，FQimg use this value to pinpoint the image in file system.
 
    Then we can get our image through the returned URL  
    https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0
    
* **image manipulation**  
 Each manipulation of the image has an corresponding command. Each command has Zero or multiple parameters. the format of the command is ```cmd``` ,which has no paramter, or ```cmd=a_b_c``` , which has three parameters a,b,c. The command attach to the rear of the URL as the standard url parameter.   
 Example:
 https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fit=200_300&grayscale  
 The above URL has two cmd 1, ```fit=200_300``` 2,```grayscale``` two commands are connected by ```&``` symbol.


## <b>command list</b>

| command name  |  command format | example  | result  |
|---|---|---|---|
|fit| fit=width_height<br>fit=width_height_filter|fit=100_300<br/>fit=100_300_box<br> fit mod cut| ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fit=150_150&raw=true)|
|fill|fill=width_height<br>fill=width_height_filter<br>fill=width_height_filter_anchor|fill=100_300<br/>fill=100_300_box<br/>fill=100_300_box_top<br>fill mod cut|![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150)|
|resize|resize=width_height<br>resize=width_height_filter|resize=100_300<br/>resize=100_300_box<br>  resize mod cut| ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?resize=150_150)|
|gamma|gamma<br>gamma=stength|gamma<br>gamma=234|![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&gamma=10)|
|contrast|contrast<br>contrast=stength|contrast<br>contrast=20<br> increase the contrast|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&contrast=120)|
|brightness|brightness<br>brightness=stength| brightness<br>brightness=0.5<br>increase the brightness |    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&brightness=38)|
|grayscale|grayscale|grayscale<br>grayscale the image |    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&grayscale)|
|invert|invert|invert<br>invert the iamge |    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&invert)|
|blur|blur<br>blur=stength|blur<br>blur=3.5<br>blur the image|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&blur=3.5)|
|sharpen|sharpen<br>sharpen=stength|sharpen<br>sharpen=3.5<br>sharpen the image|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&sharpen=65)|
|rotate90|rotate90|rotate90<br> rotate 90 degree |    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&rotate90)|
|rotate180|rotate180|rotate180<br> rotate 180 degree |    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&rotate180)|
|rotate270|rotate270|rotate270<br> rotate 270 degree |    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&rotate270)|
|flipH|flipH|flipH <br> flip horizontally|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&flipH)|
|flipV|flipV|flipV <br> flip vertically|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&flipV)|
|webp|webp<br>webp=quality| webp<br>webp=80<br/>output webp format image with 80% quality|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&webp=80)|
|jpeg|jpeg<br>jpeg=quality|jpeg<br>jpeg=80<br/>output jpeg format image with 80% quality|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&jpeg=80)|
|png|png|png <br/>output png format|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&png)|
|gif|gif<br>gif=num<br>num is number of color|png<br>png=128<br/>output png format with 128 color|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&gif=64)|
|mark|mark=id<br>mark=id_offx_offy<br>mark=id_offx_offy_offp<br>mark=id_offx_offy_offp_alpha<br>water mark <br> id is the image key set in the config file|mark=a<br>mark=a_10_10<br>mark=a_10_10_lu<br>mark=a_10_10_lu_255|    ![](https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_151&mark=a)|
