Date.prototype.Format = function (fmt) {
    var o = {
        "M+": this.getMonth() + 1,
        "d+": this.getDate(),
        "h+": this.getHours(),
        "m+": this.getMinutes(), 
        "s+": this.getSeconds(),
        "q+": Math.floor((this.getMonth() + 3) / 3),
        "S": this.getMilliseconds()
    };
    if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (var k in o)
    if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
    return fmt;
}

$(function () {
	$('.subnavbar').find ('li').each (function (i) {
		var mod = i % 3;	
		if (mod === 2) {
			$(this).addClass ('subnavbar-open-right');
		}	
	});
	initTime = new Date().getTime();
	$.getJSON("/gettime", function(out) {
		setTime(initTime, out.time);
	});
});

function setTime(initTime,serverTime) {
	ellapsedTime = new Date().getTime()-initTime;
	$('#server-time').html('当前服务器时间: <strong>'+new Date(serverTime+ellapsedTime).Format("yyyy-MM-dd hh:mm:ss")+'</strong>');
	setTimeout('setTime('+initTime+','+serverTime+');',500);
}