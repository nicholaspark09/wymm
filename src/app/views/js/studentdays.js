function DaysController(){

			var that = this;
			var days = [];
			var safekey;
			var name;
			var email;
			var phone;
			var staffRole;
			var modal;
			var table = $(document).find("#days_table");
			var moreButton = $(document).find("#viewMoreDays");
			var status = {0:"--",1:"Present",2:"Late",3:"Absent"};
			var attendances = {};

			this.Init = function(key){
				//The Safekey represents the class key
				safekey = key;
				console.log(safekey);

				Index();
				moreButton.click(function(e){
					e.preventDefault();
					Index();
				});

				$(document).on("click",".viewDay", function(e){
					e.preventDefault();
					var $that = $(this);
					var $text = $that.text();
					var ul = $that.closest('td').find(".dayTasks");
					var classdaykey = $that.attr('data-id');
					if($that.attr('data-open')=="false"){
						ul.show();
						$that.attr('data-open',"true");
					}else{
						ul.hide();
						$that.attr('data-open',"false");
					}
					var day = null;
					for(var j in days){
						if(days[j].safekey==classdaykey){
							day = days[j];
							break;
						}
					}
					if(day!=null && day.Assignments.length<1){
						ul.append("<li>Searching...</li>");
						$.get("/assignments/studentview/",{safekey:safekey,daykey:classdaykey})
							.done(function(data){
								ul.html("");
								var results = $.parseJSON(data);
								if(results['Result']=="Success")
								{
									for(var i in results['Assignments'])
									{
										var employee = results['Assignments'][i];

										AddAssignment(employee,ul);
									}	
									day.Assignments = results['Assignments'];
								}else{
									alert(results['Error']);
								}
							});
					}
				});
			}

			function AddAssignment(row,ul){
				//Print the Ass, shake it!
        		var html="";
        		if(row.pages!="")
        		{
          			var pages = $.parseJSON(row.pages);
          			html = "Pages: ";
          			for(var i in pages){
            			html+='<a href="#" class="viewPage" data-chapter="'+row.chapter+'" data-place="'+pages[i]+'"> '+pages[i]+' </a>,';
          			}
          			html = html.substring(0,html.length-1);
        		}else if(row.exam!=""){
          			html= '<a href="#" class="viewExam" data-id="'+row.exam+'"> Exam </a>';
        		}
        		//Format the due date
        		var date = row.due.substring(0,10);
        		var time = row.due.substring(11,16);
        		ul.append('<li>'+date+': '+row.name+' '+html+' <span style="weight:bold;color:red;">'+time+'</span></li>');
			}


			function Index(){
				moreButton.attr('disabled',true);
				moreButton.text("Looking...");
				moreButton.attr("class","btn btn-warning");
				$.get('/classdays/studentindex',{safekey:safekey,current: days.length})
				.done(function(data){
					var results = $.parseJSON(data);
					if(results['Result']=="Success")
					{
						var att = results['Attendances'];
						for(var i in att){
							attendances[att[i].date] = att[i];
						}
						for(var i in results['Days'])
						{
							var employee = results['Days'][i];
							employee.Assignments = [];
							days.push(employee);

							AddRow(employee);
						}
						
					}else{
						alert(results['Error']);
					}
					moreButton.attr('disabled',false);
					moreButton.text("View More");
					moreButton.attr("class","btn btn-info");
				});
			}

			function AddRow(row){
				var attendance = attendances[row.date];
				var attlink = "--";
				if(attendance!==undefined){
					attlink = status[attendance.status];
				}
					table.append('<tr><td><a href="/classdays/view/'+row.safekey+'">'+row.date.substring(0,10)+'</td><td>'+attlink+'</td><td><a href="#" class="viewDay" data-id="'+row.safekey+'" data-open="false">Assignments</a><br /><ul style="list-style-type:none;display:none;" class="dayTasks"></ul></td><td></td></tr>');
			}


		}