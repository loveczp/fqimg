<h1>go_image_server</h1>
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
<li>支持三种存储后端，本地文件，<a href="http://image.fanqiangti.net/b59fe5a3cd71bc28e39e444cd955fcb1?c=fit&w=200&h=300|c=grayscale">fastdfs</a>，<a href="https://github.com/chrislusf/seaweedfs">seaweeddfs</a></li>
</ol>

<h1>示例</h1>
原图如下
http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27

动态剪裁的例子
动态剪裁成400*400的图如下，当然高宽值可以设置成我们需要的任意值
http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=400&h=400

链式处理的例子
可以把图片裁剪成400*400后，还可以进行灰度处理。如下
<a href="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fit&w=200&h=300|c=grayscale">http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fit&w=200&h=300|c=grayscale</a>

webp处理的例子。如下
http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=webp&q=50



<h1>未来开发路线图</h1>

<ol>
<li><s>支持自定义输出格式和格式质量，比如输出jpg，png，gif。(已完成)</s></li>
<li><s>支持webp格式。webp格式压缩的图片尺寸更小。(已完成)</s></li>
<li><s>增加后端分布式存储。当前版本只支持本地文件系统存储，只能利用单台机器的存储能力。(已完成)</s></li>
<li><s>增加安全性。对于恶意攻击增加防御能力。(已完成,请见配置文件中的上传控制)</s></li>
<li>开发插件系统，使其更容已扩展</li>
</ol>


<h1>文档。</h1>
<h2>1.安装及运行</h2>
下载对应版本的可执行文件：
[windows64，linux64]('http://pan.baidu.com/s/1hr7VYle')
直接执行即可


<h2>2.图片上传</h2>
curl  --data-binary @test.jpg "http://http://image.fanqiangti.net/upload"
其中test.jpg是需要上传的文件，这个对应于http中的binary的post上传，windows平台下可以用postman来模拟。
可以得到如下结果。
{"md5":"b59fe5a3cd71bc28e39e444cd955fcb1","msg":"ok"}
md5就是文件对应的MD5码，系统也是用这个来定位上传的文件。

访问该文件方法如下
http://image.fanqiangti.net/b59fe5a3cd71bc28e39e444cd955fcb1

<h2>3.图片操作</h2>
每一个图片操作就是一个处理命令，命令有0个或者多个参数，参数包含参数名和参数值。参数是以url的参数的方式放在url尾部。即？后面就是参数。
命令名称用c表示，命令的值在下面的表中找。
例如
http://image.fanqiangti.net/b59fe5a3cd71bc28e39e444cd955fcb1?c=fit&w=200&h=300
上面的命令表示对图片进行压缩，使图片能够容纳在一个200*300的框内。命令本身参数名是c，值是fit。fit命令包含两个参数h，w分别表示高和宽。


命令之间可以通过管道链接符号“|”把多个命令连接起来实现多重操作。
例如
<a href="http://image.fanqiangti.net/b59fe5a3cd71bc28e39e444cd955fcb1?c=fit&w=200&h=300|c=grayscale">http://image.fanqiangti.net/b59fe5a3cd71bc28e39e444cd955fcb1?c=fit&w=200&h=300|c=grayscale</a>
上面表示先推图像进行fit压缩操作，然后对操作后的结果进行灰度化 处理。







<h1>命令列表</h1>

<table>
    <tr>
        <td>fit</td>
        <td>w:宽(int)
            <br/>h:高(int)</td>
        <td>c=fit&w=100&h=300<br/> fit模式裁剪</td>
        <td> <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fit&w=150&h=150" /></td>
        <td></td>
    </tr>
    <tr>
        <td>fill</td>
        <td>w:宽(int)
            <br/>h:高(int)</td>
        <td>c=fill&w=100&h=300<br/>fill模式裁剪</td>
        <td><img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150" /></td>
        <td></td>
    </tr>
    <tr>
        <td>resize</td>
        <td>w:宽(int) <br/>h:高(int) </td>
        <td> c=fill&w=100&h=300 <br/> resize模式裁剪</td>
        <td> <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=resize&w=150&h=150" /></td>
        <td></td>
    </tr>
    <tr>
        <td>gamma</td>
        <td>s:强度(float)</td>
        <td>234</td>
        <td><img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=gamma&q=10" /></td>
        <td></td>
    </tr>
    <tr>
        <td>contrast</td>
        <td>s:强度(int) </td>
        <td>c=contrast&s=20 增加对比度</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=contrast&s=120" /></td>
        <td></td>
    </tr>
    <tr>
        <td>brightness</td>
        <td>s:强度(float)</td>
        <td> c=brightness&s=0.5 增加亮度 </td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=brightness&s=38" /></td>
        <td></td>
    </tr>
    <tr>
        <td>grayscale</td>
        <td>无</td>
        <td>c=grayscale 变成灰度图 </td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=grayscale" /></td>
        <td></td>
    </tr>
    <tr>
        <td>invert</td>
        <td>无</td>
        <td>c=invert 反相 </td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=invert" /></td>
        <td></td>
    </tr>
    <tr>
        <td>blur</td>
        <td>s:强度(float)</td>
        <td>c=blur&s=3.5</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=blur&s=3.5" /></td>
        <td></td>
    </tr>
    <tr>
        <td>sharpen</td>
        <td>s:强度(float)</td>
        <td>c=sharpen&s=3.5 </td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=sharpen&s=65" /></td>
        <td></td>
    </tr>
    <tr>
        <td>rotate90</td>
        <td>无</td>
        <td>c=rotate90 正向旋转90度 </td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=rotate90" /></td>
        <td></td>
    </tr>
    <tr>
        <td>rotate180</td>
        <td>无</td>
        <td>c=rotate180正向旋转180度</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=rotate180" /></td>
        <td></td>
    </tr>
    <tr>
        <td>rotate270</td>
        <td>无</td>
        <td>c=rotate270正向旋转270度</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=rotate270" /></td>
        <td></td>
    </tr>
    <tr>
        <td>flipH</td>
        <td>无</td>
        <td>c=flipH水平翻转</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=flipH" /></td>
        <td></td>
    </tr>
    <tr>
        <td>flipV</td>
        <td>无</td>
        <td>c=flipV水平翻转</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=flipV" /></td>
        <td></td>
    </tr>
    <tr>
        <td>webp</td>
        <td>q:图片质量,可选,默认50(int) </td>
        <td> c=webp&q=80 用80%的质量输出成webp格式</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=webp&q=80" /></td>
        <td></td>
    </tr>
    <tr>
        <td>jpeg</td>
        <td>q:图片质量,可选,默认80(int) </td>
        <td>c=jpeg&q=80 用80%的质量输出成jpeg格式</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=jpeg&q=80" /></td>
        <td></td>
    </tr>
    <tr>
        <td>png</td>
        <td>无</td>
        <td>c=png 输出成png格式</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=png" /></td>
        <td></td>
    </tr>
    <tr>
        <td>gif</td>
        <td>q:颜色数量</td>
        <td>c=png&q=128 输出成128色的gif格式</td>
        <td>
            <img src="http://image.fanqiangti.net/2830dfa89daaf37b13c3421b7807df27?c=fill&w=150&h=150|c=gif&q=64" /></td>
        <td></td>
    </tr>
</table>