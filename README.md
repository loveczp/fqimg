<h1>FQimg</h1>
这是一个用go语言写的实时图片服务器
有如下特性
<ol>
<li>动态处理图片。在请求的url上加上不同的尺寸参数就可以得到不同的尺寸的图片。</li>
<li>链式图片处理。对一个图片进行多重处理，这和unix的pipeline很类似。</li>
<li>高性能。得益于go语言的并发特性，此图片服务器的内存，cpu占用都很低，能同时处理的图片数量也很可观。</li>
<li>webp图片格式支持。webp比之于jpeg，能更好的压缩图片的存储和传输体积，这点对移动应用尤为重要。</li>
<li>部署简单，只有一个可执行文件，不依赖任何外部运行库，只需将可执行文件拷贝到服务器即可执行。</li>
<li>支持多种输出格式和输出质量。当前能够支持jpeg，gif，png，bmp，webp格式，对于jpeg，gif，webp还能够支持自定义输出图片质量。</li>
<li>支持本地文件缓存，极大提高性能。</li>
<li>支持上传控制，确保只有可信IP才能上传。</li>
<li>支持跨域上传。</li>
<li>支持三种存储后端，本地文件，<a href="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fit=200_300&grayscale">fastdfs</a>，<a href="https://github.com/chrislusf/seaweedfs">seaweeddfs</a></li>
</ol>

<h1>示例</h1>
原图如下
http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27

动态剪裁的例子
动态剪裁成400*400的图如下，当然高宽值可以设置成我们需要的任意值
http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=400_400

链式处理的例子
可以把图片裁剪成400*400后，还可以进行灰度处理。如下
<a href="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fit=200_300&grayscale">http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fit=200_300&grayscale</a>

webp处理的例子。如下
http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?webp=50



<h1>未来开发路线图</h1>

<ol>
<li><s>支持自定义输出格式和格式质量，比如输出jpg，png，gif。</s>(已完成)</li>
<li><s>支持webp格式。webp格式压缩的图片尺寸更小。</s>(已完成)</li>
<li><s>增加后端分布式存储。当前版本只支持本地文件系统存储，只能利用单台机器的存储能力。</s>(已完成)</li>
<li><s>增加安全性。对于恶意攻击增加防御能力。</s>(已完成,请见配置文件中的上传控制)</li>
<li><s>增加图片水印支持。</s>(已完成,参见配置文件)</li>
<li>开发插件系统，使其更容已扩展</li>
</ol>


<h1>文档。</h1>
<h2>1.安装及运行</h2>
<p>
go get github.com/loveczp/fqimg<br/>
注意，
<ul>
<li>在windows上务必安装<a href="http://tdm-gcc.tdragon.net/download">tdm-gcc</a>否则无法编译安装通过</li>
<li>go 版本必须大于等于1.8</li>
</ul>

然后运行
fqimg -c=path/to/config/file
</p>

<h2>2.图片上传</h2>
curl -F "file=@xxxxx.jpg" "http://fqimg.com/put"
把xxxx.jpg和fqimg.com换成你对应的信息

其中test.jpg是需要上传的文件，这个对应于http中的binary的post上传，windows平台下可以用postman来模拟。
可以得到如下结果。
<br><br>
["http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27"]
<br><br>
md5就是文件对应的MD5码，系统也是用这个来定位上传的文件。

访问该文件方法如下
http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27

<h2>3.图片操作</h2>
<ul>
<li>每一个图片操作就是一个处理命令，命令有0个或者多个参数，参数包含参数名和参数值。参数是以url的参数的方式放在url尾部。即？后面就是参数。
例如
http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fit=200_300
上面的命令表示对图片进行fit压缩，使图片能够容纳在一个200*300的框内。命令是fit。对应的参数是高和宽。
</li>
<li>
命令之间可以通过连接符号“&”把多个命令连接起来实现多重操作。
例如
<a href="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fit=200_300&grayscale">http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fit=200_300&grayscale</a>
上面表示先推图像进行fit压缩操作，然后对操作后的结果进行灰度化 处理。
</li>
<ul>


<h1>开发注意事项</h1>
<ul>
<li>因为需要支持webp，所以使用了github.com/chai2010/webp 这个webp库。在<bold>windows</blod>上需要mingw的支持，请下载</li>
<li>go 版本必须大于等于1.8</li>
<ul>




<h1>命令列表</h1>

<table>
    <tr>
        <td>fit</td>
        <td> fit=width_height<br>fit=width_height_filter</td>
        <td>fit=100_300<br/>fit=100_300_box<br> fit模式裁剪</td>
        <td> <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fit=150_150" /></td>
        <td></td>
    </tr>
    <tr>
        <td>fill</td>
        <td>fill=width_height<br>fill=width_height_filter<br>fill=width_height_filter_anchor</td>
        <td>fill=100_300<br/>fill=100_300_box<br/>fill=100_300_box_top<br>fill模式裁剪</td>
        <td><img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150" /></td>
        <td></td>
    </tr>
    <tr>
        <td>resize</td>
        <td>resize=width_height<br>resize=width_height_filter</td>
        <td>resize=100_300<br/>resize=100_300_box<br>  resize模式裁剪</td>
        <td> <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?resize=150_150" /></td>
        <td></td>
    </tr>
    <tr>
        <td>gamma</td>
        <td>gamma<br>gamma=stength</td>
        <td>gamma<br>gamma=234</td>
        <td><img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&gamma=10" /></td>
        <td></td>
    </tr>
    <tr>
        <td>contrast</td>
        <td>contrast<br>contrast=stength</td>
        <td>contrast<br>contrast=20<br> 增加对比度</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&contrast=120" /></td>
        <td></td>
    </tr>
    <tr>
        <td>brightness</td>
        <td>brightness<br>brightness=stength</td>
        <td> brightness<br>brightness=0.5<br>增加亮度 </td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&brightness=38" /></td>
        <td></td>
    </tr>
    <tr>
        <td>grayscale</td>
        <td>grayscale</td>
        <td>grayscale<br>变成灰度图 </td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&grayscale" /></td>
        <td></td>
    </tr>
    <tr>
        <td>invert</td>
        <td>invert</td>
        <td>invert<br>反相 </td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&invert" /></td>
        <td></td>
    </tr>
    <tr>
        <td>blur</td>
        <td>blur<br>blur=stength</td>
        <td>blur<br>blur=3.5<br>模糊</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&blur=3.5" /></td>
        <td></td>
    </tr>
    <tr>
        <td>sharpen</td>
        <td>sharpen<br>sharpen=stength</td>
        <td>sharpen<br>sharpen=3.5<br>锐化</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&sharpen=65" /></td>
        <td></td>
    </tr>
    <tr>
        <td>rotate90</td>
        <td>rotate90</td>
        <td>rotate90 正向旋转90度 </td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&rotate90" /></td>
        <td></td>
    </tr>
    <tr>
        <td>rotate180</td>
        <td>rotate180</td>
        <td>rotate180正向旋转180度</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&rotate180" /></td>
        <td></td>
    </tr>
    <tr>
        <td>rotate270</td>
        <td>rotate270</td>
        <td>rotate270正向旋转270度</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&rotate270" /></td>
        <td></td>
    </tr>
    <tr>
        <td>flipH</td>
        <td>flipH</td>
        <td>flipH水平翻转</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&flipH" /></td>
        <td></td>
    </tr>
    <tr>
        <td>flipV</td>
        <td>flipV</td>
        <td>flipV水平翻转</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&flipV" /></td>
        <td></td>
    </tr>
    <tr>
        <td>webp</td>
        <td>webp<br>webp=quality</td>
        <td> webp<br>webp=80<br/>用80%的质量输出成webp格式</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&webp=80" /></td>
        <td></td>
    </tr>
    <tr>
        <td>jpeg</td>
        <td>jpeg<br>jpeg=quality</td>
        <td>jpeg<br>jpeg=80<br/>用80%的质量输出成jpeg格式</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&jpeg=80" /></td>
        <td></td>
    </tr>
    <tr>
        <td>png</td>
        <td>png</td>
        <td>png 输出成png格式</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&png" /></td>
        <td></td>
    </tr>
    <tr>
        <td>gif</td>
        <td>gif<br>gif=num<br>num为颜色数量</td>
        <td>png<br>png=128<br> 输出成128色的gif格式</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_150&gif=64" /></td>
        <td></td>
    </tr>
    <tr>
        <td>mark</td>
        <td>mark=mid<br>mark=mid_offx_offy<br>mark=mid_offx_offy_offp<br>mark=mid_offx_offy_offp_alpha<br>加水印 mid为配置文件中mark部分左边的key</td>
        <td>mark=a<br>mark=a_10_10<br>mark=a_10_10_lu<br>mark=a_10_10_lu_255</td>
        <td>
            <img src="http://fqimg.com/get/2830dfa89daaf37b13c3421b7807df27?fill=150_151&mark=a" /></td>
        <td></td>
    </tr>
</table>
