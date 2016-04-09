function Chapters(){
	var that = this;
	var chapters = [];
	var examkey = "";

	this.Init = function(key){
		examkey = key;
		Index();
	}

	function Index(){
		$.get('/chapters/staffindex',{current:chapters.length,safekey:examkey})
		.done(function(data){
			var results = $.parseJSON(data);
			if(results['Result']=="Success"){
				for(var i in results['Chapters']){
					var chapter = results['Chapters'][i];
					chapters.push(chapter);
				}
			}else{
				alert(results['Error']);
			}
		});
	}
}