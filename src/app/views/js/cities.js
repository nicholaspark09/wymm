function CityManager(){
					
					var that = this;
					var cities = [];
					var input = $(document).find("#cityName");
					var task = null;
					var table = $(document).find("#cities_table");
					var picked_name;
					var picked_key;

					this.Init = function(name,key){
						picked_name = name;
						picked_key = key;
						table.hide();
						input.on("keyup",function(e){
								var val = $(this).val();
								if (val.length > 2 && task == null){
									e.preventDefault();
									index(val);
								}else if (val.length==0)
								{
									task = null;
									table.hide();
								}
						});

						$(document).on("click",".clickCity", function(e){
								e.preventDefault();
								picked_name = $(this).attr("data-name");
								picked_key = $(this).attr('data-id');
								input.val(picked_name);
								table.hide();
						});

						
					};

					function index(value){
						cities = [];
						task = $.get( "http://www.tripfriends-1018.appspot.com/cities/searchindex", { current: cities.length, name: value } )
  									.done(function( data ) {
 										console.log(data);

  										var arr = $.parseJSON(data);
  										var temps = arr["Cities"];
  										table.show();
  										table.html('');
  										var html = '';
											for(var i in temps)
											{
												cities.push(temps[i]);
												html+='<tr><td class="clickCity" data-id="'+temps[i]['safekey']+'" data-name="'+temps[i]['fullname']+'">'+temps[i]['fullname']+'</td></tr>';
											}
											table.html(html);
										task = null;
  								});
					}

					this.getPicked = function(){
						return picked_key;
					}

			}