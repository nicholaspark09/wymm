function NotesManager(){
	
	var that = this;
	var notes = [];
	var controller = "";
	var action = "";
	var board;
	var saveButton;
	var currentNote = null;
	var errorLabel;
	var types = {1:"Text",2:"Handwriting",3:"Audio",4:"Video",5:"Youtube"};
	var modal;
	var moreButton;
	var Caller;
	var player = null;

	this.Init = function(callback){
		Caller = callback;
		setupModal();
		 var tag = document.createElement('script');

        tag.src = "https://www.youtube.com/iframe_api";
        var firstScriptTag = document.getElementsByTagName('script')[0];
        firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);


		$(document).on("click",".addNote", function(e){
			e.preventDefault();
			currentNote = null;
			controller = $(this).attr('data-controller');
			action = $(this).attr('data-action');
			var html = '<input type="text" placeholder="Note Title" class="form-control" id="noteName" x-webkit-speech /><br /><select id="noteType" class="form-control"><option value="0">Pick Type of Note</option>';
			for(var i in types){
				html+='<option value="'+i+'">'+types[i]+'</option>';
			}
			html+='</select>';
			$(document).find("#noteHeader").html(html);
			$(document).find("#noteName").val("");
			modal.modal('show');
			$(document).find("#noteType").val(1);
			board.html('<textarea id="noteBody" style="width:100%;height:250px;box-shadow:0 0 5px #ccc;" x-webkit-speech></textarea>');

		});

		$(document).on("change","#noteType", function(e){
			e.preventDefault();
			var type = $(this).val();
			if(currentNote==null){
				if(type==1){
					board.html('<textarea id="noteBody" style="width:100%;height:250px;box-shadow:0 0 5px #ccc;" x-webkit-speech></textarea>');
				}else if(type==5){
					board.html('<input type="text" id="noteBody" placeholder="Youtube ID" class="form-control" /><br /><br /><div id="player"></div>');
				}
			}else{
				if(type==1){
					board.html('<textarea id="noteBody" style="width:100%;height:250px;box-shadow:0 0 5px #ccc;" x-webkit-speech>'+currentNote.prettybody+'</textarea>');
				}else if(type==5){
					board.html('<input type="text" id="noteBody" value="'+currentNote.prettybody+'" class="form-control" /><br /><br /><div id="player"></div>');
					width=320;
					height=240;
					player = new YT.Player('player', {
                            height: height,
                            width: width,
                            videoId: currentNote.prettybody,
                          });
				}
			}
		});

		$(document).on("blur","#noteBody", function(e){
			var val = $(document).find("#noteType").val();
			console.log(val);
			if(val==5){
	
					width=320;
					height=240;
					player = new YT.Player('player', {
                            height: height,
                            width: width,
                            videoId: $(document).find("#noteBody").val(),
                          });
			}
		});

		$(document).on("click","#saveNote", function(e){
			if(currentNote == null){
				var note = {name:$(document).find("#noteName").val(),body:$(document).find("#noteBody").val(),type:$(document).find("#noteType").val()};
				Add(note);
			}else{
				currentNote.name = $(document).find("#noteName").val();
				currentNote.body = $(document).find("#noteBody").val();
				currentNote.type = $(document).find("#noteType").val();
				Edit();
			}
			
		});

		$(document).on("click",".viewNote", function(e){
			
			e.preventDefault();
			var id = $(this).attr('data-id');
			View(id);
		});

		$(document).on("click",".deleteNote", function(e){
			
			e.preventDefault();
			var id = $(this).attr('data-id');
			var $that = $(this);
			$that.attr('disabled',true);
			$that.text("...");
			var url = '/notes/delete/'+id;
			$.get(url)
			.done(function(data){
				console.log(data);
				var results = $.parseJSON(data);
				if(results['Result']=="Success"){
					$that.closest('li').remove();
				}else{
					alert(results['Error']);
					$that.attr('disabled',false);
					$that.text("x");
				}
			});
		});
	}

	function setupModal(){
		var html = '<div class="modal fade" id="notesCreator" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">\
  					<div class="modal-dialog" role="document">\
    					<div class="modal-content">\
     					 	<div class="modal-header" id="noteHeader">\
     					 	<input type="text" placeholder="Note Title" class="form-control" id="noteName" x-webkit-speech />\
      						</div>\
      						<div class="modal-body" id="noteBoard" style="height:290px;overflow:auto;">\
      						</div>\
      					<div class="modal-footer">\
      						<div id="noteErrorLabel" style="color:red;"></div>\
        					<button type="button" class="btn btn-info" id="saveNote">Save</button>\
        					<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>\
      					</div>\
    					</div>\
  					</div>\
					</div>';
		$(html).appendTo('body');
		board = $(document).find("#noteBoard");
		saveButton = $(document).find("#saveNote");
		errorLabel = $(document).find("#noteErrorLabel");
		modal = $(document).find("#notesCreator");
	}

	function Add(note){
		if(note.body==""){
			alert("Please insert a title");
		}else{
			errorLabel.text("Saving...");
			$.ajax({
				type: 'POST',
				url: '/notes/add',
				data:{
					'controller':controller,
					'action':action,
					'body':note.body,
					'name':note.name,
					'type':note.type
				},
				success: function(data){
					console.log(data);
					var results= $.parseJSON(data);
					if(results['Result']=="Success"){	
						note.prettybody = note.body;
						notes.push(note);
						View(note.safekey);
						Caller.AddNotes([note]);
						errorLabel.text("Saved");
					}else{
						errorLabel.text(results['Error']);
					}
				},
				error: function(){

				}
			});
		}
	}

	function Edit(){
		errorLabel.text("Saving....");
			$.ajax({
				type: 'POST',
				url: '/notes/edit',
				data:{
					'safekey':currentNote.safekey,
					'body':currentNote.body,
					'name':currentNote.name,
					'type':currentNote.type
				},
				success: function(data){
					console.log(data);
					var results= $.parseJSON(data);
					if(results['Result']=="Success"){	
						errorLabel.text("Saved");
					}else{
						errorLabel.text(results['Error']);
					}
				},
				error: function(){
					errorLabel.text('No connection');
				}
			});
	}

	function View(safekey){
		currentNote = null;
		for(var i in notes){
			if(notes[i].safekey==safekey){
				currentNote = notes[i];
				break;
			}
		} 
		if(currentNote==null){
		$.get('/notes/webview/',{safekey:safekey})
		.done(function(data){
			var results = $.parseJSON(data);
			if(results['Result']=="Success"){
				var note = results['Note'];
				notes.push(note);
				currentNote = note;
				Display();
			}else{
				alert(results['Error']);
			}
		});
		}else
			Display();
	}

	function Display(){
		$(document).find("#noteName").val(currentNote.name);
		var html = '<input type="text" placeholder="Note Title" class="form-control" id="noteName" x-webkit-speech /><br /><select id="noteType" class="form-control"><option value="0">Pick Type of Note</option>';
			for(var i in types){
				html+='<option value="'+i+'">'+types[i]+'</option>';
			}
			html+='</select>';
		$(document).find("#noteHeader").html(html);
		$(document).find("#noteName").val(currentNote.name);
		$(document).find("#noteType").val(currentNote.type);
		var body = "";
		if(currentNote.type==1){
			
			$(document).find("#noteType").val(currentNote.type);
			board.html('<textarea id="noteBody" style="width:100%;height:250px;box-shadow:0 0 5px #ccc;" x-webkit-speech>'+currentNote.prettybody+'</textarea>');
			
			modal.modal('show');
		}else if(currentNote.type==5){
			$(document).find("#noteType").val(currentNote.type);
			board.html('<input type="text" id="noteBody" value="'+currentNote.prettybody+'" class="form-control" /><br /><br /><div id="player"></div>');
			width=320;
					height=240;
					player = new YT.Player('player', {
                            height: height,
                            width: width,
                            videoId: currentNote.prettybody,
                          });
			modal.modal('show');
		}
		$(document).find("#noteType").val(currentNote.type);
	}

	//Callback must include function called AddNotes
	this.Index = function(controller,action,current){
		$.get('/notes/webindex',{current:current,controller:controller,action:action})
		.done(function(data){
			console.log(data);
			var results = $.parseJSON(data);
			if(results['Result']=="Success"){

				Caller.AddNotes(results['Notes']);

			}else{
				Caller.ShowError(results['Error']);
			}
		});
	}

}