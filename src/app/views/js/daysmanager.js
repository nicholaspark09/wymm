function ClassdaysManager(){
	
	var that = this;
	var days = [];
	var table;
	var searchDate;
	var moreButton;
	var searchButton;

	this.Init = function(table_id){
		table = $(document).find("#"+table_id);
		//Setup the calendar
		jQuery('#agendaDate').datetimepicker({format:'Y-m-d'});
		searchDate = $(document).find("#agendaDate");
		moreButton = $(document).find("#viewMoreDays");
		searchButton = $(document).find("#searchAgenda");
		current = new Date();
        var yyyy = current.getFullYear().toString();
        var mm = (current.getMonth()+1).toString(); // getMonth() is zero-based
        var dd  = current.getDate().toString();
        date = yyyy+"-"+(mm[1]?mm:"0"+mm[0])+"-"+(dd[1]?dd:"0"+dd[0]); // padding
        searchDate.val(date);
        Index();

        searchButton.click(function(e){
        	e.preventDefault();
        	days = [];
        	table.html('');
        	Index();
        });
	}

	function Index(){
		if(searchDate.val()=="")
		{
			alert("Please enter a date");
		}else{
			moreButton.attr('disabled',true);
			moreButton.text("Loading...");
			moreButton.attr("class","btn btn-warning");
			searchButton.attr('disabled',true);
			var url = "/classdays/index";
			$.get(url,{date:searchDate.val(),current:days.length})
			.done(function(data){
				console.log(data);
				var results = $.parseJSON(data);
				if(results['Result']=="Success"){
					for(var i in results['Classdays']){
						var day = results['Classdays'][i];
						AddRow(day);
					}
					if(results['Classdays'].length<50){
						moreButton.text("No more");
						moreButton.attr('class','btn btn-info');
					}else{
						moreButton.attr('disabled',false);
						moreButton.attr('class','btn btn-default');
						moreButton.text("View More");
					}
				}else
				{
					alert(results['Error']);
					moreButton.attr('disabled',false);
					moreButton.attr('class','btn btn-info');
					moreButton.text("View More");
				}
				searchButton.attr('disabled',false);
			});
		}
	}

	function AddRow(row){
		table.append('<tr><td><a href="/classdays/view/'+row.safekey+'">'+row.beg.substring(11,16)+' - '+row.end.substring(11,16)+'</a></td><td><a href="/classes/view/'+row.coursedaykey+'">'+row.name+'</a></td></tr>');
	}

}