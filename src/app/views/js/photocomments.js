function PhotoComments(){
		var that = this;
		var comments = [];
		var main;
		that.cont = "";
		that.action = "";
		var table;
		var saveButton;
		var body;
		var errorLabel;
		var profiles = [];
		var order = 0;
		var isMore = true;
		var finding = {};
		var currentComment;

		this.Init = function(div,cont,action)
		{
			that.cont = cont;
			that.action = action;
			main = $(document).find("#"+div);
			that.Reset(that.action);



			//Like it
					$(document).on("click",".likePhoto"+that.cont+"It", function(e){
						e.preventDefault();
						var $that = $(this);
						var controller = "comments";
						var act = $that.attr('data-id');
						$that.attr('disabled',true);
						$that.text("Liking...");
						$.ajax({
								type: 'POST',
								url: '/likes/add',
								data:{
									'cont':controller,
									'action':act
								},
								success: function(data){
									var results = $.parseJSON(data);
									if(results['Result']=="Success"){
										$that.text("Liked");
										$that.attr('class','unlikePhoto'+that.cont+'It');
									}else{
										$that.text("Like");
									}
									$that.attr('disabled',false);
								},
								error: function(){
									showMessage(true,"No connection");
									$that.text("Like");
									$that.attr('disabled',false);
								}
						});
					});

					//Unlike it
					$(document).on("click",".unlikePhoto"+that.cont+"It", function(e){
						e.preventDefault();
						var $that = $(this);
						var controller = "comments";
						var act = $that.attr('data-id');
						$that.attr('disabled',true);
						$that.text("Unliking...");
						$.ajax({
								type: 'POST',
								url: '/likes/delete',
								data:{
									'cont':controller,
									'action':act
								},
								success: function(data){
									console.log(data);
									var results = $.parseJSON(data);
									if(results['Result']=="Success"){
										$that.text("Like");
										$that.attr('class','likePhoto'+that.cont+'It');
									}else{
										$that.text("Liked");
									}
									$that.attr('disabled',false);
								},
								error: function(){
									showMessage(true,"No connection");
									$that.text("Liked");
									$that.attr('disabled',false);
								}
						});
					});

					$(document).on("click",".deletePhoto"+that.cont+"Comment", function(e){
							e.preventDefault();
							var $that = $(this);
							var id = $that.attr('data-id');
							$that.attr('disabled',true);
							if(confirm("Are you sure?"))
							{
								$that.text("Deleting...");
								$.get( "/comments/delete/"+id)
											.done(function(data){
												var results = $.parseJSON(data);
												if(results['Result']=="Success")
													$that.closest('tr').remove();
												else{
													$that.text("Delete");
													$that.attr('disabled',false);
													errorLabel.text(results['Error']);
												}
								});
							}
					});

					$(document).on("click",".editPhoto"+that.cont+"Comment", function(e){
						e.preventDefault();
						var $that = $(this);
						var p = $that.closest('tr').find(".photoCommentBody");
						currentComment = $that.attr('data-id');
						p.html('<textarea class="current'+that.cont+'CommentEdit">'+p.text()+'</textarea>');
					});

					$(document).on("keyup",".current"+that.cont+"CommentEdit", function(e){
						var $that = $(this);

						if(e.keyCode == 13)
						{
						var commentBody = $that.val();
						$.ajax({
								type: 'POST',
								url: '/comments/edit/'+currentComment,
								data:{
									'body':commentBody
								},
								success: function(data){
									console.log(data);
									var results = $.parseJSON(data);
									if(results['Result']=="Success")
									{
										showMessage(false,"Comment Saved");
										$that.closest('p').html(commentBody);
									}else{
										showMessage(true,results['Error']);
									}
						
								},
								error: function(){
									showMessage(true,"No connection");
				
								}
						});
						}
					});

					$(document).on("click","#loadMorePhotos", function(e){
						e.preventDefault();
						var tr = $(this).closest('tr');
						tr.remove();
						index();

					});	

					$(document).on("keyup","#photoCommentBody", function(e){
						if(e.keyCode==13)
						{
							e.preventDefault();
							var text = $(this).val();
							if(text==""){
								showMessage(true,"Please write in some text");
							}else{
								saveButton.attr('disabled',true);
								add(text);
							}	
						}
					});
		}

		function add(message){
	
						showMessage(false,"Saving...");
						$.ajax({
							type: 'POST',
							url: '/comments/add/',
							data:{
								'body':message,
								'cont':that.cont,
								'action':that.action
							},
							success: function(data){
								var results = $.parseJSON(data);
		
								if(results['Result'] == "Success"){
									showMessage(false,"Saved");
									var commentkey = results['Key'];
									var safeuser = results['UserKey'];
									table.prepend('<tr><td style="background:white;width:90%;height:auto;border:1px solid #ccc;"><h4><a href="/profiles/view/'+safeuser+'" class="profileName" data-found="false" data-id="'+safeuser+'" id="profile'+safeuser+'">Me</a></h4><br /><p style="text-indent:20px;margin-bottom:0;" class="photoCommentBody" id="commentBody'+commentkey+'">'+message+'</p><a href="#" class="deletePhoto'+that.cont+'Comment" style="color:red;float:right;" data-id="'+commentkey+'">Delete</a><a href="#" class="editPhoto'+that.cont+'Comment" style="margin-right:10px;float:right;" data-id="'+commentkey+'" >Edit</a><a href="#" class="likePhoto'+that.cont+'It" data-id="'+commentkey+'" style="margin-right:10px;float:right;">Like</a></td></tr>');

								}else{
									showMessage(true,results['Error']);
								}
								saveButton.attr('disabled',false);
								task = null;
								$(document).find("#photoCommentBody").val("");
							},
							error: function(){
								showMessage(true,"No Connection");
								task = null;
								saveButton.attr('disabled',false);
							}
						});
					
				}

		function index(){
					showMessage(false,"Loading comments...");
					$.get( "/comments/index/",{current:comments.length,cont:that.cont,action:that.action,order:order.val()})
											.done(function(data){
													console.log(data);
													var results = $.parseJSON(data);
													if(results['Result']=="Success")
													{

														for(var i in results['Comments']){
															var comment = results['Comments'][i];
															comments.push(comment);
															var link = '<a href="#" class="likePhoto'+that.cont+'It" data-id="'+comment.safekey+'" style="margin-right:10px;float:right;">Like ('+comment.likes+')</a>';
															if(comment.ilike==true){
																link = '<a href="#" class="unlikePhoto'+that.cont+'It" data-id="'+comment.safekey+'" style="margin-right:10px;float:right;">Liked ('+comment.likes+')</a>'
															}
															table.append('<tr><td style="background:white;width:90%;height:auto;border:1px solid #ccc;"><h4><a href="/profiles/view/'+comment.userkey+'" class="profileName" data-found="false" data-id="'+comment.userkey+'" id="profile'+comment.userkey+'">'+comment.name+'</a></h4><br /><p style="text-indent:20px;" class="photoCommentBody" id="commentBody'+comment.safekey+'">'+comment.body+'</p><a href="#" class="deletePhoto'+that.cont+'Comment" style="color:red;float:right;" data-id="'+comment.safekey+'">Delete</a><a href="#" class="editPhoto'+that.cont+'Comment" style="margin-right:10px;float:right;" data-id="'+comment.safekey+'">Edit</a>'+link+'</td></tr>');

														}
														console.log("The length is "+results['Comments'].length);
														if(results['Comments'].length<15)
															isMore = false;
														else{
															isMore = true;
															
															table.append('<tr class="loadPhotoRow"><td><a href="#" class="btn btn-info" id="loadMorePhotos">Load More </a></td></tr>');
														}
														display();
														showMessage(false,"");
													}else
													{
														showMessage(true,results['Error']);
												
													}

													
					});
				}

		function display(){
			$(document).find(".profileName").each(function(e){
						var $that =$(this);
						var found = $that.attr('data-found');
						var id = $that.attr('data-id');
						var localFind = false;
						for(var j in profiles){
							
							if(profiles[j].user == id){
								localFind = true;
								$that.html('<img src="'+profiles[j].pic+'" style="width:60px;" /> '+profiles[j].name);
								break;
							}
						}
						if(found=="false" && !localFind && finding[id]===undefined){
							finding[id] = true;
							$.get( "/profiles/grabmini/"+id)
											.done(function(data){
								
												var results = $.parseJSON(data);
												if(results['Result']=="Success")
												{
													var profile = results['Profile'];
													if(results['Pic']===undefined)
														profile.pic="";
													else
														profile.pic = "data:image/png;base64, "+results['Pic'];
													profiles.push(profile);
													$(document).find(".profileName").each(function(){
														var id = $that.attr('data-id');
														if(id == profile.user){
															$(this).html('<img src="'+profile.pic+'" style="width:60px;" /> '+profile.name);
											
														}
												
														$(this).attr('data-found','true');
													});
												}
											});
						}

					});
		}

		this.Reset = function(action)
		{
			that.action = action;
			comments = [];
			main.html('<h2>Comments</h2><div id="photoErrorLabel"></div><textarea  style="width:250px;height:100px;" id="photoCommentBody"></textarea><a href="#" id="savePhoto'+that.cont+'Comment" class="btn btn-success">Save</a><br /><table><tr><td><select id="photoCommentOrder"><option value="0">Newest-Oldest</option><option value="1">Oldest-Newest</option><option value="2">Most Liked</option><option value="3">Least Liked</option></select></td></tr><tbody id="photo_comments_table"></tbody></table>');
			table = $(document).find("#photo_comments_table");
			saveButton = $(document).find("#savePhoto"+that.cont+"Comment");
			body = $(document).find("#photoCommentBody");
			errorLabel = $(document).find("#photoErrorLabel");
			order = $(document).find("#photoCommentOrder");



			saveButton.click(function(e){
				e.preventDefault();
						var text = body.val();
						if(text==""){
							showMessage(true,"Please write in some text");
						}else{
							saveButton.attr('disabled',true);
							add(text);
						}
			});
			index();
		}

		function showMessage(error, message)
				{
					if(error){
						errorLabel.css({'color':'red'});
						errorLabel.text(message);
					}else{
						errorLabel.css({'color':'green'});
						errorLabel.text(message);
					}
				}
}