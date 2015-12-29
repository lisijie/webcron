
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
	$('#server-time').html('当前服务器时间: <strong>'+new Date(serverTime+ellapsedTime).toLocaleString()+'</strong>');
	setTimeout('setTime('+initTime+','+serverTime+');',500);
}