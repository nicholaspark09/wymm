function TaskCreator(){
	
	var that = this;
	var Books;
	var classkey = "";
	var classday = "";
	var board;
	var modal;
	var createButton;
	var nextStep;
	var name;
	var start;
	var end;
  var MediaMaker;
  var taskTypes = {1:"Book",2:"Instruction",3:"Multimedia",4:"Youtube",5:"Exam",6:"Otherday"};
  var selectedType = 1;


	this.Init = function(class_key,media){
		classkey = class_key;
    if(media!=null)
      MediaMaker = media;
		SetupModal();
		$(document).on("click",".addSingleTask", function(e){
              e.preventDefault();
              createButton.show();
              nextStep.hide();
              classday = $(this).attr('data-id');
              modal.modal('show');
              board.html('<ul style="list-style-type:none;"><a href="#" id="addmultimediaTask" class="newTaskType" data-type="3"><li style="float:left;margin:10px;padding:10px;box-shadow:0 0 5px #ccc;">Multimedia Task</li></a><a href="#" id="singleTask" class="newTaskType" data-type="2"><li style="float:left;margin:10px;padding:10px;box-shadow:0 0 5px #ccc;">Instructions</li></a><a href="#" id="bookTask" class="newTaskType" data-type="1"><li style="float:left;margin:10px;padding:10px;box-shadow:0 0 5px #ccc;">Book Task</li></a><a href="#" class="courseTask" data-type="6" class="newTaskType"><li style="float:left;margin:10px;padding:10px;box-shadow:0 0 5px #ccc;">Course Catalog</li></a><a href="#" id="examTask" data-type="5" class="newTaskType"><li style="float:left;margin:10px;padding:10px;box-shadow:0 0 5px #ccc;">Exam Task</li></a><a href="#" id="youtubeTask" data-type="4" class="newTaskType"><li style="float:left;margin:10px;padding:10px;box-shadow:0 0 5px #ccc;">Youtube Task</li></a></ul>');
          });

		$(document).on("click","#addmultimediaTask", function(e){
			e.preventDefault();
			createButton.hide();
			nextStep.show();
			board.html("<h4>Create the task first then click next</h4>");
		});

    $(document).on("click",".newTaskType", function(e){
      e.preventDefault();
      selectedType = $(this).attr('data-type');
      if(selectedType==2){
        board.html("Enter instructions above:");
      }else if(selectedType==4){
        board.html('Please insert Youtube id in the title');
      }
    });

    $(document).on("blur","#taskName", function(e){
      if(selectedType==4){
        var id = $(this).val();
        board.html('<iframe width="350" height="290" src="https://www.youtube.com/embed/'+id+'" frameborder="0" allowfullscreen></iframe>');
      }
    });

		$(document).on("click","#nextStep", function(e){
			e.preventDefault();
			nextStep.hide();
			console.log(name.val());
			console.log(end.val());
			Add(name.val(),"","","",start.val(),end.val(),classday,function(tasks){
				if(tasks.length==1){
					console.log("New task "+tasks[0].name);
          MediaMaker.NewVision("assignments",tasks[0].safekey)
				}
			});
			board.html("Loading Multimedia adder");
		});

    createButton.click(function(e){
      e.preventDefault();
      Add(name.val(),"","","", start.val(),end.val(),classday,function(tasks){
        console.log("Succeeded");
        console.log(tasks.length);
      });
    });
		 //Delete task
          $(document).on("click",".deleteAssignment", function(e){
            e.preventDefault();
            if(confirm("Are you sure?"))
            {
              var $that = $(this);
              var id = $that.attr('data-id');
              $that.text("Deleting...");
              var url = '/assignments/staffdelete/'+id;
              $.get(url,{parentkey:classkey})
              .done(function(data){
                console.log(data);
                var results= $.parseJSON(data);
                if(results['Result']=="Success")
                {
                  $that.closest('tr').remove();
                }else{
                  alert(results['Error']);
                }
              });
            }
          });
	}

	function SetupModal(){
		html = '<div class="modal fade" id="taskboard" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">\
  					<div class="modal-dialog" role="document">\
    					<div class="modal-content">\
     					 	<div class="modal-header">\
   							<input type="text" placeholder="Task Description" class="form-control" id="taskName" />\
          						<ul id="taskPages" style="list-style-type:none;">\
          						</ul>\
      						</div>\
      						<div class="modal-body" id="taskBoard" style="height:250px;overflow:auto;">\
      						</div>\
      					<div class="modal-footer">\
        					<input type="text" class="form-control" id="taskStart" placeholder="Start date/Time"/><br />\
        					<input type="text" class="form-control" id="taskEnd" placeholder="End Date / Due TIME!"/>\
        					<button type="button" class="btn btn-info" id="nextStep" style="display:none;">Next</button>\
        					<button type="button" class="btn btn-info" id="createTask">Create</button>\
        					<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>\
      					</div>\
    					</div>\
  					</div>\
					</div>';
		$(html).appendTo('body');
		board = $(document).find("#taskBoard");
          modal = $(document).find("#taskboard");
          createButton = $(document).find("#createTask");
          nextStep = $(document).find("#nextStep");
          jQuery("#taskStart").datetimepicker({format:'Y-m-d H:i'});
          jQuery("#taskEnd").datetimepicker({format:'Y-m-d H:i'});
          name = $(document).find('#taskName');
          start = $(document).find('#taskStart');
		  end = $(document).find('#taskEnd');
	}

	function Add(name,book,chapter,pages,start,due,classday,callback){
	if(start == "")
	{
		alert("Please include a date");
	}else{
          $.ajax({
                  type: 'POST',
                  url: '/assignments/classadd',
                  data:{
                    'name':name,
                    'bookkey':book,
                    'chapter':chapter,
                    'pages':pages,
                    'start':start,
                    'due':due,
                    'classkey':classkey,
                    'classday':classday,
                    'type':taskTypes[selectedType]
                  },
                  success: function(data){
                    console.log(data);
                    var results = $.parseJSON(data);
                    if(results['Result']=="Success")
                    {
                      alert("Saved");
                      if(callback!=null){
                      	callback(results['Assignments']);
                      }
                    }else
                      alert(results['Error']);
                  },
                  error: function(){
                    alert("No connection");
                  }
              });
		}
        }

}