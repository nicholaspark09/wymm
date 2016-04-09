function AttendanceManager(){
	
	var that = this;
	var attends = [];
	var table;
	var moreButton;
	var classdayKey;
	var statuses = {0:"-",1:": )",2:": (",3:"X"};
	var ids = [];
	var callback;
	var today = "";

	this.Init = function(table_id,safekey,func){
		table = $(document).find("#"+table_id);
		moreButton = $(document).find("#viewMoreAttends");
		calculateToday();
		classdayKey = safekey;
		callback = func;
		Index();
		moreButton.click(function(e){
			e.preventDefault();
			Index();
		});

		$(document).on("click",".saveDate", function(e){
			e.preventDefault();
			var $that = $(this);
			var id = $that.attr('data-id');
			var status = $that.closest('tr').find('.setAtt').val();
			var datestamp = $that.closest('tr').find('.dateStamp').val();
			if(datestamp == "")
			{
				alert("Please fill in the time");
			}else{
				$that.text("Saving...");
				$that.attr('disabled',true);
				$.ajax({
					type: 'POST',
					url: '/attendances/edit',
					data:{
						'safekey':id,
						'status':status,
						'date':datestamp
					},
					success: function(data){
						var results = $.parseJSON(data);
						if(results['Result']=="Success"){
							$that.css({'color':'green'});
						}else{
							$that.css({'color':'red'});
							alert(results['Error']);
						}
						$that.text("Save");
						$that.attr('disabled',false);
					},error: function(){
						alert("No connection");
						$that.text("Save");
						$that.attr('disabled',false);
					}
				});
			}
		});
	}

	function calculateToday(){
		var current = new Date();
        var yyyy = current.getFullYear().toString();
        var mm = (current.getMonth()+1).toString(); // getMonth() is zero-based
        var dd  = current.getDate().toString();
        today = date = yyyy+"-"+(mm[1]?mm:"0"+mm[0])+"-"+(dd[1]?dd:"0"+dd[0])+" "+current.getHours()+":"+current.getMinutes();
	}


	function Index(){
		moreButton.text("Loading...");
		moreButton.attr('disabled',true);
		var url = "/attendances/webindex/"+classdayKey;
		$.get(url,{current:attends.length})
		.done(function(data){
			var results = $.parseJSON(data);
			if(results['Result']=="Success"){
				for(var i in results['Attendances'])
				{
					var attend = results['Attendances'][i];
					attends.push(attend);
					AddRow(attend);
				}
				callback(ids);
				$(document).find(".setAtt").each(function(){
					var status = $(this).attr('data-status');
					var i = parseInt(status);
					$(this).val(i);
				});
			}else
			{
				alert(results['Error']);
				moreButton.text("More");
				moreButton.attr('disabled',false);
			}
		});
	}

	function AddRow(row){
		var status = '<select class="setAtt" data-id="'+row.safekey+'" data-status="'+row.status+'">';
		for(var i in statuses){
			status+='<option value="'+i+'">'+statuses[i]+'</option>';
		}
		var stampDate = row.arrival.substring(0,10)+" "+row.arrival.substring(11,16);
		if(stampDate == "0001-01-01 00:00"){
			stampDate = today;
		}
		var id = "dateId"+row.safekey;
		table.append('<tr><td><a href="/profiles/view/'+row.student+'">'+row.name+'</a></td><td>'+status+'</td><td><input type="text" id="'+id+'" class="dateStamp" value="'+stampDate+'" /></td><td><a href="#" class="saveDate" data-id="'+row.safekey+'">Save</a></td><td><a href="#" class="deleteAttend" data-id="'+row.safekey+'">Delete</a></td></tr>');
		ids.push(id);
	}

}