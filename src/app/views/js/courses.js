function CoursesManager(){
	
	var that = this;
	var courses = [];
	var campusKey = "";
	var table;
	var catalog = {};
	var classManager = null;

	this.Init = function(key){

		campusKey = key;
		table = $(document).find("#courses_table");
		Index();

		$(document).on("click",".addCatalog", function(e){
			e.preventDefault();
			var $that = $(this);
			Add($that);
		});

		$(document).on("click",".deleteCourse", function(e){
			e.preventDefault();
			var $that = $(this);
			if(confirm("Delete?"))
			{
				$that.text("Deleting...");
				var url = '/courselists/delete/'+$that.attr('data-id');
				$.get(url)
				.done(function(data){
					var results = $.parseJSON(data);
					if(results['Result'] == "Success")
					{
						$that.closest('tr').remove();
					}else{
						$that.text("Delete");
						alert(results['Error']);
					}
				});
			}
			
		});
	}

	function Add(button){
		button.attr('disabled',true);
		button.text("Saving...");
		$.ajax({
			type: 'POST',
			url: '/courselists/add',
			data:{
				'campuskey':campusKey,
				'coursekey':button.attr('data-id')
			},
			success: function(data){
				console.log(data);
				var results = $.parseJSON(data);
				if(results['Result']=="Success")
				{
					button.attr('disabled',true);
					button.text("Added!");
					Index();
				}else{
					button.attr('disabled',false);
					button.text("Add");
				}
			},
			error: function(){
				alert("No connection");
				button.attr('disabled',false);
				button.text("Add");
			}
		});
	}

	function Index(){
		var url = '/courselists/staffindex/'+campusKey;
		$.get(url,{current:courses.length})
		.done(function(data){
			console.log(data);
			var results= $.parseJSON(data);
			if(results['Result']=="Success")
			{
				for(var i in results['Courses']){
					catalog[results['Courses'][i].safekey] = results['Courses'][i].name;
				}
				for(var i in results['Courselists'])
				{
					var courselist = results['Courselists'][i];
					courses.push(courselist);
					AddRow(courselist);
				}
	
			}else
				alert(results['Error']);
		});
	}

	function AddRow(courselist){
		table.append('<tr><td><a href="/courselists/adminview/'+courselist.safekey+'">'+catalog[courselist.course]+'</a></td><td><a href="#" class="deleteCourse" data-id="'+courselist.safekey+'">Delete</a></td></tr>');
	}

}