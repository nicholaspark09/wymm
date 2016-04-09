/*
	Messages Controller


		In order to send messages to someone, please insert a link with class ".sendMessage"
*/
function MessagesManager(){
	var that = this;
	var safekey = "";
	var otherkey = "";
	var threads = [];
	var messages = [];
	var safekeys = [];
	var modal;
	var sendButton;
	var errorLabel;
	var isSending = false;
	var body;

	this.Init = function(){

		Setup();

		$(document).on("click",".sendMessage", function(e){
			e.preventDefault();
			otherkey = $(this).attr('data-id');
			modal.modal('show');
			sendButton = $(this);
			sendButton.attr('disabled',true);
			safekey = "";
			View();
		});

		$(document).on("keyup","#messageInput",function(e){
			if(e.keyCode==13){
				e.preventDefault();
				var data = $(this).val();
				Add(data);
			}
		});

		$(document).on("click","#closeNewModal",function(e){
			e.preventDefault();
			sendButton.attr('disabled',false);
		});
	}

	function Add(input){
		if(isSending==false)
		{
			isSending = true;
			errorLabel.text("Sending...");
			$.ajax({
				type: 'POST',
				url: '/messagethreads/add',
				data:{
					'safekey':safekey,
					'body':input
				},
				success: function(data){
					var results = $.parseJSON(data);
					if(results['Result']=="Success"){
						errorLabel.text("");
						message = {name:"Me",body:input,mine:true,seen:0};
						messages.push(message);
						AddRow(message);
					}else{
						errorLabel.text(results['Error']);
					}
				},
				error: function(){
					errorLabel.text("No connection");
					isSending = false;
				}
			});	
	    }
	}

	function View(){
		errorLabel.text("Loading...");
		$.ajax({
			type: 'POST',
			url: '/messagethreads/view',
			data:{
				'other':otherkey
			},
			success: function(data){
				var results = $.parseJSON(data);
				if(results['Result']=="Success"){
					safekey = results['SafeKey'];
					$(document).find("#messageTitle").text(results['Name']);
					errorLabel.text('');
					Index();
				}else{
					if(results['Name']!==undefined){
						$(document).find("#messageTitle").text(results['Name']);
					}
					errorLabel.text(results['Error']);
				}
				sendButton.attr('disabled',false);
			},
			error: function(){
				errorLabel.text("No connection");
			}
		});
	}



	function Index(){
		if(safekey != "")
		{
			//get the current message thread
			errorLabel.text("Loading....");
			$.get('/messagethreads/indexmessages',{current:messages.length,safekey:safekey})
			.done(function(data){
				var results = $.parseJSON(data);
				errorLabel.text('');
				if(results['Result']=="Success"){
					var temp = results['Messages'];
					temp.reverse();
					for(var i in temp){
						var message = temp[i];
						message.mine = true;
						if(message.User == otherkey){
							message.mine = false;
						}
						messages.push(message);
						AddRow(message);
					}
				}else{
					errorLabel.text(results['Error']);
				}
			});
		}
	}

	function AddRow(message){
		if(message.mine == true){
			body.append('<div style="margin-top:5px;background:blue;color:white;padding:10px;margin-left:20px;float:right;text-align:right;clear:both;">'+message.body+'</div>');
		}else{
			body.append('<div style="margin-top:5px;background:#ccc;color:black;padding:10px;margin-left:20px;float:left;tex-align:left;clear:both;">'+message.body+'</div>');
		}
		console.log(message);
	}

	function Setup(){
		var html = '<div class="modal fade" id="newMessageModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">\
		  <div class="modal-dialog" role="document">\
		    <div class="modal-content">\
		      <div class="modal-header">\
		          <h1 id="messageTitle"></h1>\
		          <div id="messageErrorLabel" style="color:red;font-weight:bold;">\
		      	</div>\
		      </div>\
		      <div class="modal-body" id="messageBody" style="height:250px;overflow:scroll;">\
		      </div>\
		      <div class="modal-footer">\
		      	<input type="text" class="form-control" id="messageInput" placeholder="Write your message" />\
		        <button type="button" id="closeNewModal" class="btn btn-default" data-dismiss="modal">Close</button>\
		      </div>\
		    </div>\
		  </div>\
		</div>';
		$(document.body).append(html);
		modal = $(document).find("#newMessageModal");
		errorLabel = $(document).find("#messageErrorLabel");
		body = $(document).find("#messageBody");
	}
}