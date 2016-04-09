function CoursebookManager(){
	
	var that =this;
	var classkey = "";
	var books = [];
	var table;
	var board;
	var currentBook = null;
	var pageView;
	var chaptersButton;
	var chaptersList;
	var pagesList;
	var currentPage = null;
	var currentChapter = null;
	
	this.Init = function(class_key){
		classkey = class_key;


		//Setup the book viewer
		SetupDisplay();

		$(document).on("click",".viewCoursebook", function(e){
			e.preventDefault();
			var $that = $(this);
			var book = {name:$that.text(),safekey:$that.attr('data-id'),Bookplace:null};
			currentBook = null;
			for(var i in books){
				if(books[i].safekey == book.safekey){
					currentBook= books[i];
					break;
				}
			}
			if(currentBook!=null)
			{
				if(currentBook.Chapters.length==0)
					Index();
				else{
					DisplayBook();
				}
			}else{
				books.push(book);
				book.Chapters = [];
				currentBook = book;
				Index();
			}
		});

		$(document).on("change","#pagesList", function(e){
			e.preventDefault();
			ShowPage();
		});

		$(document).on("click","#closeBookView", function(e){
			e.preventDefault();
			board.hide();
		});


		$(document).on("change",".chaptersList", function(e){
			SetupChapter();
		});
	}

	function SetupChapter(){
		currentChapter = null;
		currentBook.Bookplace.chapterkey = chaptersList.val();
		for(var i in currentBook.Chapters[i]){
			if(currentBook.Chapters[i].safekey == currentBook.Bookplace.chapterkey){
					currentChapter = currentBook.Chapters[i];
					break;
			}
		}
		if(currentChapter!=null){
			if(currentChapter.Pages.length<1){
				var url = "/coursebooks/classpages/";
				$.get(url,{chapter:currentChapter.safekey,classkey:classkey})
				.done(function(data){
					console.log(data);
					var results = $.parseJSON(data);
					pagesList.html('');
					if(results['Result']=="Success"){
						for(var i in results['Pages']){
							var page = results['Pages'][i];
							currentChapter.Pages.push(page);
							pagesList.append('<option value="'+page.place+'">'+page.place+'</option>');
						}
						pagesList.val(currentBook.Bookplace.place);
						ShowPage();
					}else{
						alert(results['Error']);
					}
				});
			}else{
				pagesList.html('');
				for(var i in currentChapter.Pages){
					pagesList.append('<option value="'+currentChapter.Pages[i].place+'">'+currentChapter.Pages[i].place+'</option>');
				}
				pagesList.val(currentBook.Bookplace.place);
				ShowPage();
			}
		}
	}

	function SetupDisplay(){
		board = $('<div id="board" style="width:100%;height:100%;position:fixed;top:0;left:0;margin:0;z-index:9995;background:white;display:none;"><div id="bookPageView" style="width:100%;height:95%;top:0;left:0;text-align:center;margin:0;overflow:auto;"></div><ul id="bookControlView" style="width:100%;height:5%;position:fixed;bottom:0;left:0;list-style-type:none;border-top:1px solid #ccc;"><li style="line-height:25px;width:20%;float:left;text-align:center;"><select id="chaptersList">Chapter</select></li><li style="line-height:25px;width:20%;float:left;text-align:center;"><select id="pagesList">Pages</select></li><a href="#" id="closeBookView" class="btn btn-danger"><li style="line-height:25px;width:20%;float:left;text-align:center;">x</li></a></ul></div>').appendTo('body');
		pageView = $(document).find("#bookPageView");
		chaptersButton = $(document).find("#viewChapters");
		chaptersList = $(document).find("#chaptersList");
		pagesList = $(document).find("#pagesList");
	}

	function DisplayBook(){
		pageView.html("Loading...");
		board.show();
		chaptersList.html('');
		pagesList.html('<li>Loading...</li>');
		for(var i in currentBook.Chapters){
			chaptersList.append('<option value="'+currentBook.Chapters[i].safekey+'">'+currentBook.Chapters[i].name+'</option>');
		}
		currentChapter = null;
		if(currentBook.Bookplace!=null){
			console.log("Did you find it");
			for(var i in currentBook.Chapters){
				if(currentBook.Chapters[i].safekey == currentBook.Bookplace.chapterkey){
					currentChapter = currentBook.Chapters[i];
					break;
				}
			}
			if(currentChapter.Pages.length<1){
				var url = "/coursebooks/classpages/";
				$.get(url,{chapter:currentChapter.safekey,classkey:classkey})
				.done(function(data){
					console.log(data);
					var results = $.parseJSON(data);
					pagesList.html('');
					if(results['Result']=="Success"){
						for(var i in results['Pages']){
							var page = results['Pages'][i];
							currentChapter.Pages.push(page);
							pagesList.append('<option value="'+page.place+'">'+page.place+'</option>');
						}
						pagesList.val(currentBook.Bookplace.place);
						ShowPage();
					}else{
						alert(results['Error']);
					}
				});
			}else{
				pagesList.html('');
				for(var i in currentChapter.Pages){
					pagesList.append('<option value="'+currentChapter.Pages[i].place+'">'+currentChapter.Pages[i].place+'</option>');
				}
				pagesList.val(currentBook.Bookplace.place);
				ShowPage();
			}
		}else{
			currentChapter = currentBook.Chapters[0];
			if(currentChapter.Pages.length<1){
				var url = "/coursebooks/classpages/";
				$.get(url,{chapter:currentChapter.safekey,classkey:classkey})
				.done(function(data){
					console.log(data);
					var results = $.parseJSON(data);
					pagesList.html('');
					if(results['Result']=="Success"){
						for(var i in results['Pages']){
							var page = results['Pages'][i];
							currentChapter.Pages.push(page);
							pagesList.append('<option value="'+page.place+'">'+page.place+'</option>');
						}
						ShowPage();
					}else{
						alert(results['Error']);
					}
				});
			}else{
				pagesList.html('');
				for(var i in currentChapter.Pages){
					pagesList.append('<option value="'+currentChapter.Pages[i].place+'">'+currentChapter.Pages[i].place+'</option>');
				}
				ShowPage();
			}
		}
		
	}

	function ShowPage(){
		var val = pagesList.val();
		for(var i in currentChapter.Pages){
			if(currentChapter.Pages[i].place == val){
				currentPage = currentChapter.Pages[i];
				break;
			}
		}
		if(currentPage.type==1)
			pageView.html(currentPage.prettybody);
		else if(currentPage.type==2)
		{
			width = $(document).width();
			height = $(document).height()*.7;
			pageView.html('<iframe width="'+width+'" height="'+height+'" src="https://www.youtube.com/embed/'+currentPage.video+'" frameborder="0" allowfullscreen></iframe>');
		}
		$.ajax({
			type:'POST',
			url: '/bookplaces/add',
			data:{
				'chapter':currentChapter.safekey,
				'place':currentPage.place
			},
			success: function(data){
				console.log(data);
				var results = $.parseJSON(data);
				if(results['Result']=="Success"){
					var bookplace = results['Bookplace'];
					currentBook.Bookplace = bookplace;
				}else	
					console.log(results['Error']);
			},
			error: function(){
				console.log("No connection");
			}
		})
	}


	function Index(){
		var url = "/coursebooks/classchapters/"+currentBook.safekey;
		$.get(url,{classkey:classkey})
		.done(function(data){
			console.log(data);
			var results= $.parseJSON(data);
			if(results['Result']=="Success"){
				for(var i in results['Chapters']){
					var chapter = results['Chapters'][i];
					chapter.Pages = [];
					currentBook.Chapters.push(chapter);
				}
				if(results['Bookplace']!=""){
					currentBook.Bookplace = results['Bookplace'];
				}
				DisplayBook();
			}
			else{
				alert(results['Error']);
			}
		});
	}


}