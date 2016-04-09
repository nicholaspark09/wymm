function AssignmentsManager(){
	

	/*
		The permissions key is the key against which the Role Model checks 
			If you have permission to access whatever permission key you give, then you have access to add,update,modify or delete the task
			i.e.
				If you're adding an assignment to a Courseday, the Permissions Key would be the CourseKey
																The Parent Key would be the CoursedayKey
	*/

	var that =this;
	var permissions_key;
	var parent_key;
	var modal;
	var board;
	var taskTypes = {1:"Book",2:"Instruction",3:"Multimedia",4:"Youtube",5:"Exam",6:"Sound Cloud",7:"Manual Exam",8:"Event"};
	var selectedType = 1;
	var Books = null;
	var catalog = [];
	var pickedBook;
	var pickedChapter;
	var MediaMaker;
	var createButton;
	var currentExam="";
	var inClass = false;
	var baseType = "";
	var points = 0;
	var required = false;

	this.Init = function(safekey,media){
		permissions_key = safekey;
		MediaMaker = media;
		SetupModal();
		$(document).on("click",".addTask", function(e){
			e.preventDefault();
			parent_key = $(this).attr('data-id');
			RefreshModal();
		});

		$(document).on("click",".newTaskOfType", function(e){
			e.preventDefault();
			selectedType = $(this).attr('data-id');
			if(selectedType==1){
				if(Books != null && catalog.length==0){
					catalog = Books.getBooks();
				}
				var ul = '<ul style="list-style-type:none;"">';
				for(var i in catalog){
					catalog[i].Exams = [];
					console.log(catalog[i].name);
						ul+='<a href="#" class="pickBook" data-id="'+catalog[i].safekey+'"><li style="padding:10px;border:1px solid #ccc;border-radius:10px;margin:10px;text-align:center;">'+catalog[i].name+'</li></a>';;
				}
				ul+='</ul>';
				board.html(ul);
				$(document).find("#taskPages").show();
			}else if(selectedType==2){
				board.html("<h2>Please write your instructions in the title up top :) </h2>");
			}else if(selectedType==3){
				board.html('<h2>Add the Title, pick the requirement, and click Create</h2>');
			}else if(selectedType==4){
				board.html("<h2>Please put the Youtube Id in the title</h2>");
			}else if(selectedType==5){
				var html = 'Enter Title of Exam and <br /><br /><b>Choose Type:</b> \
				<select id="taskExamType" class="form-control">\
              	<option value="1">Entrance</option>\
             	 <option value="2">Pop Quiz</option>\
              	<option value="3">Quiz</option>\
             	 <option value="4">Vocab Quiz</option>\
              	<option value="5">Chapter Review</option>\
              	<option value="6">Semester Exam</option>\
              	<option value="7">Mid-Term Exam</option>\
              	<option value="8">Final Exam</option>\
              	<option value="9">Evaluation Exam</option>\
            	</select>';
            	board.html(html);
			}else if (selectedType == 6){
				board.html("<h2>Please insert sound cloud link above</h2>");
			}else if (selectedType == 7){
				board.html("<h2>Please write the title of the exam above. Remember, this is a physical exam not a digital one.</h2>");
			}
		});

		$(document).on("click","#restartTasks", function(e){
			e.preventDefault();
			$(document).find("#taskPages").hide();
			currentExam = "";
			RefreshModal();
		});

		$(document).on("blur","#taskName", function(e){
			if(selectedType==4){
				var id = $(this).val();
				board.html('<iframe width="350" height="290" src="https://www.youtube.com/embed/'+id+'" frameborder="0" allowfullscreen></iframe>');
			}
		});

		$(document).on("click","#chooseChapters",function(e){
			e.preventDefault();
			selectedType = 1;
			currentExam = "";
			ShowChapters();
		});

		$(document).on("click","#chooseBookExam",function(e){
			e.preventDefault();
			selectedType = 7;
			if(pickedBook.Exams.length==0)
				GetExams();
			else
				ShowExams();
			//Seven is a temporary book exam, this will revert back later
		});

		$(document).on("change","#rowThreeExam",function(e){
      			var val = $(this).val();
      			if(val!="0")
      			{
        			currentExam = val;
      			}else{
        			currentExam = "";
      			}
    	});

		//You picked the book!
		$(document).on("click",".pickBook", function(e){
						e.preventDefault();
						var $that = $(this);
						pickedBook = null;
						//Clear the modal first
					
						var id = $that.attr('data-id');
						for(var i in catalog){
							if(catalog[i].safekey == id){
								pickedBook = catalog[i];
								break;
							}
						}
							var html= '<table class="table table-striped"><tr><th>'+pickedBook.name+'</th><th><a href="#" id="chooseChapters">Chapters</a><br /><a href="#" id="chooseBookExam">Exams</a></th></tr><tbody id="boardTable"></tbody></table>';
							board.html(html);
					});	

			$(document).on("click",".pickChapter", function(e){
				e.preventDefault();
				pickedChapter = null;
				var $that = $(this);
				var picked = $that.attr('data-picked');
				if(picked=="false"){
						var ul = $that.closest('td').find(".chapterPages");
						var td = $that.closest('td').next('td');
						var id = $that.attr('data-id');
						for(var i in pickedBook.Chapters){
							if(pickedBook.Chapters[i].safekey == id){
								pickedChapter = pickedBook.Chapters[i];
								break;
							}
						}
						if(pickedChapter!=null && pickedChapter.Pages.length>0){
							for(var i in pickedChapter.Pages){
								var page= pickedChapter.Pages[i];
										ul.append('<li><input type="checkbox" class="addPage" data-id="'+page.safekey+'" data-place="'+page.place+'" /><a href="#" class="pickPage" data-id="'+page.safekey+'">'+page.place+'</a></li>');
							}

						}else if(pickedChapter.Pages.length<1){
							$.get("/coursebooks/staffpages/",{course:safekey,chapter:id})
							.done(function(data){
								console.log(data);
								var results = $.parseJSON(data);
								ul.append('<li><a href="#" class="pickAllPages" data-picked="false">Pick All</a></li>');
								if(results['Result']=="Success")
								{	
									for(var i in results['Pages']){
										var page= results['Pages'][i];
										pickedChapter.Pages.push(page);
										ul.append('<li><input type="checkbox" class="addPage" data-id="'+page.safekey+'" data-place="'+page.place+'" /><a href="#" class="pickPage" data-id="'+page.safekey+'">'+page.place+'</a></li>');
									}
								}else{
									alert(results['Error']);
								}
							});
						}
					ul.show();
					$that.attr('data-picked',"true");
				}else{
					var ul = $that.closest('td').find(".chapterPages");
					ul.hide();
					$that.attr('data-picked',"false");
				}
			});

			$(document).on("click",".pickAllPages", function(e){
						var $that = $(this);
						var ul = $(document).find("#taskPages");
						ul.html('');
						if($that.attr('data-picked')=="false")
						{
							$(this).closest('ul').find("input").each(function(){
								$(this).prop('checked',true);
								var id = $(this).attr('data-id');
								var place= $(this).attr('data-place');
								ul.append('<a href="#" class="pickedPage" data-id="'+id+'" data-place="'+place+'"><li style="float:left;margin:5px;">'+place+'</li></a>');
							});
							$that.attr('data-picked',"true");
						}else{
							$(this).closest('ul').find("input").each(function(){
								$(this).prop('checked',false);
							});
							$that.attr('data-picked',"false");
						}
					});

			$(document).on("click",".pickPage", function(e){
						e.preventDefault();
						var id = $(this).attr('data-id');
						var $that = $(this);
						pickedPage = null;
						for(var i in pickedChapter.Pages){
							if(pickedChapter.Pages[i].safekey == id){
								pickedPage = pickedChapter.Pages[i];
								break;
							}
						}
						$that.closest('td').next('td').html('<div style="height:150px;overflow:auto;">'+pickedPage.prettybody+'</div><br /> Page '+$that.text());
					});

			$(document).on("change",".addPage", function(e){
						e.preventDefault();
						var $that = $(this);
						var ul = $(document).find("#taskPages");
						var id = $(this).attr('data-id');
						var place= $(this).attr('data-place');
						if($that.is(":checked"))
						{
							ul.append('<a href="#" class="pickedPage" data-id="'+id+'" data-place="'+place+'"><li style="float:left;margin:5px;">'+place+'</li></a>');
						}else{
							ul.find(".pickedPage").each(function(){
								if($(this).attr('data-id')==id){
									$(this).remove();
									return;
								}
							});
						}
			});

			$(document).on("click","#createTasks", function(e){
				e.preventDefault();

				var name = $(document).find("#taskName").val();
				if($(document).find("#taskRequired").is(":checked"))
					required =true;
				else
					required = false;
				if($(document).find("#taskInClass").is(":checked"))
					inClass = true;
				else
					inClass = false;
				points = $(document).find("#taskPoints").val();
				baseType = $(document).find("#taskBaseType").val();
				if(selectedType==1){
					var pages = [];
					$(document).find(".pickedPage").each(function(){
								var id = $(this).attr('data-place');
								pages.push(id);
							}).promise().done(function(){

								Add(name,pickedBook.safekey,pickedChapter.safekey,JSON.stringify(pages));
							});
					
				}else if(selectedType==2){
					Add(name,"","","");
				}else if(selectedType==3){
					Add(name,"","","");
				}else if(selectedType==4){
					Add(name,"","","");
				}else if(selectedType==5){
					//Create tehe exam first
					createButton.attr('disabled',true);
					createButton.text("Creating Exam...");
					 $.ajax({
          					type: 'POST',
          					url: '/exams/staffadd/',
          					data:{
          						'parentkey':permissions_key,
          						'childkey':parent_key,
          						'name':name,
          						'description':'',
          						'type':$(document).find("#taskExamType").val()
          					},
          					success: function(data){
          						var results= $.parseJSON(data);
								createButton.text("Create");
          						if(results['Result']=="Success"){
          							
          							currentExam = results['Key'];
          							console.log(currentExam);
          							Add(name,"","","");
          						}else{
          							createButton.attr('disabled',false);
          							alert(results['Error']);
          						}
          						
          					},
          					error: function(){
          						alert("No Connection");
          					}

          				});
					
				}else if(selectedType==6){
					Add(name,"","","");
				}else if(selectedType==7){
					Add(name,"","","");
				}
			});
	}

	function RefreshModal(){
		inClass = false;
		baseType = "";
		required = false;
		points = 0;
		var html = '<ul style="list-style-type:none;margin:10px;display:block;">';
			for(var i in taskTypes){
				html+='<a href="#" class="newTaskOfType" data-id="'+i+'"><li style="padding:10px;border:1px solid #ccc;border-radius:10px;margin:10px;text-align:center;">'+taskTypes[i]+'</li></a>';
			}
			html+='</ul>';
			board.html(html);
	}

	function SetupModal(){
		var html = '<div class="modal fade" id="taskAdder" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">\
  				<div class="modal-dialog" role="document">\
    				<div class="modal-content">\
      					<div class="modal-header">\
   							<h4>Add Task</h4>\
        						<input type="text" id="taskName" class="form-control" placeholder="Title / Directions" /><br />\
        						<ul style="list-style-type:none;" id="taskPages">\
        						</ul>\
      					</div>\
      					<div class="modal-body" id="taskAdderBoard" style="height:auto">\
      					</div>\
      					<div class="modal-footer">\
      						<ul style="list-style-type:none;" id="bookChooser">\
      							</ul>\
      						<b>Points</b><input type="number" id="taskPoints" /><br />\
      						<b>In Class</b><input type="checkbox" id="taskInClass" /><br />\
      						<b>Category</b><select id="taskBaseType"><option value="">Category</option></select><br />\
      						<b>Required?</b> <input type="checkbox" id="taskRequired" /><br />\
      						<button type="button" class="btn btn-info" id="restartTasks">Restart</button>\
      						<button type="button" class="btn btn-info" id="createTasks">Create</button>\
        					<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>\
      					</div>\
    				</div>\
  				</div>\
			</div>';
			$(html).appendTo('body');
			modal = $(document).find("#taskAdder");
			board = $(document).find("#taskAdderBoard");
			createButton = $(document).find("#createTasks");
	}

	function Add(name,book,chapter,pages){
		createButton.attr('disabled',true);
		createButton.text("Saving...");
		console.log(currentExam);

		

			$.ajax({
                  type: 'POST',
                  url: '/assignments/staffadd',
                  data:{
                    'name':name,
                    'coursekey':permissions_key,
                    'bookkey':book,
                    'courseday':parent_key,
                    'chapter':chapter,
                    'pages':pages,
                    'exam':currentExam,
                    'type':taskTypes[selectedType],
                    'points':points,
                    'inclass':inClass,
                    'required':required,
                    'basetype':baseType
                  },
                  success: function(data){
                    var results = $.parseJSON(data);
                    if(results['Result']=="Success"){
                      RefreshModal();
                      modal.modal('hide');
                      if(selectedType==3){
                      	MediaMaker.NewVision("assignments",results['Key']);
                      }else if(selectedType==5){
                      	var key = currentExam;
                      	currentExam="";
                      	window.location = '/exams/staffedit/'+key;
                      }
                    }else{
                      alert(results['Error']);
                    }
                    createButton.text("Create");
                    createButton.attr('disabled',false);
                  },
                  error: function(){
                    alert("No connection");
                    createButton.text("Create");
                    createButton.attr('disabled',false);
                  }
                });
   }

   function ShowChapters(){
   	selectedType = 1;
   	boardTable = $(document).find("#boardTable");
   	boardTable.html('');
   		if(pickedBook!=null && pickedBook.Chapters.length>0){
							for(var i in pickedBook.Chapters){
										var chapter = pickedBook.Chapters[i];
										boardTable.append('<tr><td><a href="#" class="pickChapter" data-id="'+chapter.safekey+'" data-picked="false">'+chapter.name+'</a><br /><ul style="list-style-type:none;" class="chapterPages"></ul></td></tr>');
							}
						}else if(pickedBook.Chapters.length==0){
							var url = "/coursebooks/staffchapters/"+pickedBook.safekey;
							console.log(url);
							$.get(url)
							.done(function(data){
								var results = $.parseJSON(data);
								if(results['Result']=="Success")
								{	
									for(var i in results['Chapters']){
										var chapter = results['Chapters'][i];
										chapter.Pages = [];
										pickedBook.Chapters.push(chapter);
										boardTable.append('<tr><td><a href="#" class="pickChapter" data-id="'+chapter.safekey+'" data-picked="false">'+chapter.name+'</a><br /><ul style="list-style-type:none;" class="chapterPages"></ul></td><td></td></tr>');
									}
								}else{
									alert(results['Error']);
								}
							});
						}
   }

   function GetExams(){
   	var url = "/coursebooks/staffexams/"+pickedBook.safekey;
      $.get(url)
      .done(function(data){
        var results= $.parseJSON(data);
        if(results['Result']=="Success"){
          for(var i in results['Exams']){
            var chapter = results['Exams'][i];
            pickedBook.Exams.push(chapter);
          }
        }
        else{
          alert(results['Error']);
        }
        ShowExams();
      });
   }

   function ShowExams(){
   	var boardTable = $(document).find("#boardTable");
   	  var td = $(document).find("#rowThree");
      var html = '<tr><td><select id="rowThreeExam" style="form-control"><option value="0">Pick Exam</option>';
      for(var i in pickedBook.Exams){
        html+='<option value="'+pickedBook.Exams[i].safekey+'">'+pickedBook.Exams[i].name+'</option>';
      }
      html+='</select></td></tr>';
      boardTable.html(html);
   }

   this.setBooks = function(books){
		Books = books;
		catalog = Books.getBooks();
		for(var i in catalog){
			catalog[i].Exams = [];
		}
	}

	this.setTypes = function(types){
		console.log(types);
		var select =$(document).find("#taskBaseType");
		for(var i in types){
			select.append('<option value="'+types[i].name+'">'+types[i].name+'</option>');
		}	
	}

}