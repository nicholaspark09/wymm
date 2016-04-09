function CManager(){
				var that = this;
				var comments = [];
				var container;

				var task = null;
				that.cont = "";
				that.action = "";
				var ul;
				var profiles = [];
				var finding = {};
				var currentComment = "";
				var order = 0;
				var isMore = true;

				this.Init = function(commentDiv,controller,action){
					that.cont = controller;
					that.action = action;
					container = $(document).find("#"+commentDiv);
					setupDivs();

					//ul = $(document).find("#commentUL");
					
					index();

					$(document).on("click",".delete"+that.cont+"Comment", function(e){
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
													$that.closest('li').remove();
												else{
													$that.text("Delete");
													$that.attr('disabled',false);
													errorLabel.text(results['Error']);
												}
								});
							}
					});

					$(document).on("click",".edit"+that.cont+"Comment", function(e){
						e.preventDefault();
						var $that = $(this);
						var p = $that.closest('li').find(".commentBody");
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

					$(document).on("change",".comment"+that.cont+"Order", function(e){
						e.preventDefault();
						comments = [];
						ul.html('');
						index();
					});

					//Like it
					$(document).on("click",".like"+that.cont+"It", function(e){
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
										$that.attr('class','unlike'+that.cont+'It');
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
					$(document).on("click",".unlike"+that.cont+"It", function(e){
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
										$that.attr('class','likeIt');
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

					$(document).on("click","#loadMore", function(e){
							e.preventDefault();
							if(isMore){
								$(this).remove();
								index();

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
							console.log(profiles[j]);
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

				function index(){
					showMessage(false,"Loading comments...");
					order = $(document).find(".comment"+that.cont+"Order").val();
					$.get( "/comments/index/",{current:comments.length,cont:that.cont,action:that.action,order:order})
											.done(function(data){
													console.log(data);
													var results = $.parseJSON(data);
													if(results['Result']=="Success")
													{

														for(var i in results['Comments']){
															var comment = results['Comments'][i];
															comments.push(comment);
															var link = '<a href="#" class="like'+that.cont+'It" data-id="'+comment.safekey+'" style="margin-right:10px;float:right;">Like ('+comment.likes+')</a>';
															if(comment.ilike==true){
																link = '<a href="#" class="unlike'+that.cont+'It" data-id="'+comment.safekey+'" style="margin-right:10px;float:right;">Liked ('+comment.likes+')</a>'
															}
															ul.append('<br /><li style="background:white;width:90%;height:auto;padding:10px;border:1px solid #ccc;border-radius:10px;padding-bottom:20px;"><h4><a href="/profiles/view/'+comment.userkey+'" class="profileName" data-found="false" data-id="'+comment.userkey+'" id="profile'+comment.userkey+'">'+comment.name+'</a></h4><br /><p style="text-indent:20px;" class="commentBody" id="commentBody'+comment.safekey+'">'+comment.body+'</p><a href="#" class="delete'+that.cont+'Comment" style="color:red;float:right;" data-id="'+comment.safekey+'">Delete</a><a href="#" class="edit'+that.cont+'Comment" style="margin-right:10px;float:right;" data-id="'+comment.safekey+'">Edit</a>'+link+'</li>');

														}
														console.log("The length is "+results['Comments'].length);
														if(results['Comments'].length<15)
															isMore = false;
														else{
															isMore = true;
															
															ul.append('<a href="#" class="btn btn-info" id="load'+that.cont+'More">Load More </a>');
														}
														display();
														showMessage(false,"");
													}else
													{
														showMessage(true,results['Error']);
												
													}

													
					});
				}

				function add(message){
					if(checkTask){
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
									ul.prepend('<br /><li style="background:white;width:90%;height:auto;padding:10px;border:1px solid #ccc;border-radius:10px;padding-bottom:20px;"><h4><a href="/profiles/view/'+safeuser+'" class="profileName" data-found="false" data-id="'+safeuser+'" id="profile'+safeuser+'">Me</a></h4><br /><p style="text-indent:20px;" class="commentBody" id="commentBody'+commentkey+'">'+message+'</p><a href="#" class="delete'+that.cont+'Comment" style="color:red;float:right;" data-id="'+commentkey+'">Delete</a><a href="#" class="edit'+that.cont+'Comment" style="margin-right:10px;float:right;" data-id="'+commentkey+'" >Edit</a><a href="#" class="like'+that.cont+'It" data-id="'+commentkey+'" style="margin-right:10px;float:right;">Like</a></li>');

								}else{
									showMessage(true,results['Error']);
								}
								saveButton.attr('disabled',false);
								task = null;
							},
							error: function(){
								showMessage(true,"No Connection");
								task = null;
								saveButton.attr('disabled',false);
							}
						});
					}
				}

				//Checks for running async tasks, if there are you can cancel
				//returns true if you're ready to go!
				function checkTask(){
					if(task!= null){
						if(confirm("Break the old task to run your new one?"))
						{
							task.abort();
							task = null;
							return true;
						}else
							return false;
					}
					return true;
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

				function setupDivs(){
					container.html('');
					container.append('<div class="commentErrorLabel"></div><textarea  style="width:250px;height:100px;" class="commentBody"></textarea><br /><a href="#" class="saveComment">Save</a><br /><select class="comment'+that.cont+'Order"><option value="0">Newest-Oldest</option><option value="1">Oldest-Newest</option><option value="2">Most Liked</option><option value="3">Least Liked</option></select>');
					container.append('<ul style="list-style-type:none;" class="commentUL"></ul>');
					ul = container.find(".commentUL");
					saveButton = container.find(".saveComment");
					body = container.find(".commentBody");
					errorLabel = container.find(".commentErrorLabel");
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
				}

				that.resetSettings = function(newAction){
					console.log("What?");
					that.action = newAction;
					comments = [];
					setupDivs();
					index();
				}

			}