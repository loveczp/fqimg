# FQimg 

先来体验一下把  
http://loveczp.github.io/fqimg/

这是一个用go语言写的实时图片服务器
有如下特性
* 动态处理图片。在请求的url上加上不同的尺寸参数就可以得到不同的尺寸的图片。
* 链式图片处理。对一个图片进行多重处理，这和unix的pipeline很类似。
* 高性能。得益于go语言的并发特性，此图片服务器的内存，cpu占用都很低，能同时处理的图片数量也很可观。
* webp图片格式支持。webp比之于jpeg，能更好的压缩图片的存储和传输体积，这点对移动应用尤为重要。
* 部署简单，只有一个可执行文件，不依赖任何外部运行库，只需将可执行文件拷贝到服务器即可执行。
* 支持多种输出格式和输出质量。当前能够支持jpeg，gif，png，bmp，webp格式，对于jpeg，gif，webp还能够支持自定义输出图片质量。
* 支持本地文件缓存，极大提高性能。
* 支持上传控制，确保只有可信IP才能上传。
* 支持跨域上传。
* 支持http 和 https
* 支持两种存储后端，本地文件， [seaweeddfs](https://github.com/chrislusf/seaweedfs)

# 示例 
* **原图如下**  
https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0

* **动态剪裁的例子**  
动态剪裁成400*400的图如下，当然高宽值可以设置成我们需要的任意值  
https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=400_400

* **链式处理的例子**  
可以把图片裁剪成400*400后，还可以进行灰度处理。如下  
https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fit=200_300&grayscale

* **webp处理的例子。如下**  
https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?webp=50





# 文档。 
* 安装及运行 
    ```bash
    go get github.com/loveczp/fqimg
    ```
    注意:
    * 在windows上务必安装<a href="http://tdm-gcc.tdragon.net/download">tdm-gcc</a>否则无法编译安装通过
    * go 版本必须大于等于1.8

    然后运行
    ```bash
    fqimg -c=path/to/config/file
    ```

* 图片上传 
   ```bash
    curl -F "file=@xxxxx.jpg" "http://fqimg.com/put"
   ```
    把xxxx.jpg和fqimg.com换成你对应的信息可以得到如下结果。  
    ["https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0"]  
    md5就是文件对应的MD5码，系统也是用这个来定位上传的文件。

    访问该文件方法如下
    https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0



* 3.图片操作 
    * 每一个图片操作就是一个处理命令，命令有0个或者多个参数，参数包含参数名和参数值。参数是以url的参数的方式放在url尾部。即？后面就是参数。  
    例如  
    https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fit=200_300  
    上面的命令表示对图片进行fit压缩，使图片能够容纳在一个200*300的框内。命令是fit。对应的参数是高和宽。

    * 命令之间可以通过连接符号“&”把多个命令连接起来实现多重操作。
    例如  
    https://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fit=200_300&grayscale
    上面表示先推图像进行fit压缩操作，然后对操作后的结果进行灰度化 处理。




# 命令列表

| 命令名称  |  命令格式 | 例子  | 结果  |
|---|---|---|---|
|fit|fit=width_height<br>fit=width_height_filter|        fit=100_300<br/>fit=100_300_box<br> fit模式裁剪|![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fit=150_150) |
|       fill|        fill=width_height<br>fill=width_height_filter<br>fill=width_height_filter_anchor|        fill=100_300<br/>fill=100_300_box<br/>fill=100_300_box_top<br>fill模式裁剪|        ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150) |
|       resize|        resize=width_height<br>resize=width_height_filter|        resize=100_300<br/>resize=100_300_box<br>  resize模式裁剪|![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?resize=150_150) |
|       gamma|        gamma<br>gamma=stength|        gamma<br>gamma=234|        ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&gamma=10) |
|       contrast|        contrast<br>contrast=stength|        contrast<br>contrast=20<br> 增加对比度|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&contrast=120) |
|       brightness|        brightness<br>brightness=stength|         brightness<br>brightness=0.5<br>增加亮度 |    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&brightness=38) |
|       grayscale|        grayscale|        grayscale<br>变成灰度图 |    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&grayscale) |
|       invert|        invert|        invert<br>反相 |    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&invert) |
|       blur|        blur<br>blur=stength|        blur<br>blur=3.5<br>模糊|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&blur=3.5) |
|       sharpen|        sharpen<br>sharpen=stength|        sharpen<br>sharpen=3.5<br>锐化|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&sharpen=65) |
|       rotate90|        rotate90|        rotate90 正向旋转90度 |    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&rotate90) |
|       rotate180|        rotate180|        rotate180正向旋转180度|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&rotate180) |
|       rotate270|        rotate270|        rotate270正向旋转270度|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&rotate270) |
|       flipH|        flipH|        flipH水平翻转|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&flipH) |
|       flipV|        flipV|        flipV水平翻转|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&flipV) |
|       webp|        webp<br>webp=quality|         webp<br>webp=80<br/>用80%的质量输出成webp格式|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&webp=80) |
|       jpeg|        jpeg<br>jpeg=quality|        jpeg<br>jpeg=80<br/>用80%的质量输出成jpeg格式|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&jpeg=80) |
|       png|        png|        png 输出成png格式|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&png) |
|       gif|        gif<br>gif=num<br>num为颜色数量|        png<br>png=128<br> 输出成128色的gif格式|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_150&gif=64) |
|       mark|        mark=mid<br>mark=mid_offx_offy<br>mark=mid_offx_offy_offp<br>mark=mid_offx_offy_offp_alpha<br>加水印 mid为配置文件中mark部分左边的key|        mark=a<br>mark=a_10_10<br>mark=a_10_10_lu<br>mark=a_10_10_lu_255|    ![](http://fqimg.com/get/0657cae447e8c88f44c65b7e5f73cfe0?fill=150_151&mark=a) |
