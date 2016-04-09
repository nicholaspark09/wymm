function StaffExams(){
	
	var that = this;
	var parentkey = "";
	var childkey = "";
	var exams = [];
	var moreButton = $(document).find("#viewMoreExams");
	var table;

	this.Init = function(parent,child,tableDiv){
		parentkey = parent;
		childkey = child;
		table = $(document).find("#"+tableDiv);
		$(document).on("click","#createExam", function(e){
			e.preventDefault();
			var name = $(document).find("#examName").val();
			var description = $(document).find("#examDescription").val();
			var thetype = $(document).find("#examType").val();
			if(name=="")
				alert("Please fill in the Exam title");
			else{
				Add(name,description,thetype);
			}
		});

		$(document).on("click",".deleteExam", function(e){
			e.preventDefault();
			if(confirm("Delete?"))
			{
				var id = $(this).attr('data-id');
				DeleteExam(id,$(this));
			}
		});
		Index();
	}

	function Add(name,desc,thetype){
		var button = $(document).find("#createExam");
		button.text("Creating...");
		button.attr('class','btn btn-warning');
		button.attr('disabled',true);
		$.ajax({
			type: 'POST',
			url: '/exams/staffadd',
			data:{
				'parentkey':parentkey,
				'childkey':childkey,
				'name':name,
				'description':desc,
				'type':thetype
			},
			success: function(data){
				var results = $.parseJSON(data);
				if(results['Result']=="Success"){
					alert("Saved!");
					var exam = {name:name,description:desc,safekey:results['Key'],parentkey:parentkey};
					AddRow(exam);
				}else
					alert(results['Error']);
				button.text("Create");
				button.attr('class','btn btn-info');
				button.attr('disabled',false);
			},
			error: function(){
				alert("No connection");
				button.text("Create");
				button.attr('class','btn btn-info');
				button.attr('disabled',false);
			}
		});
	}

	//Get the child key (current Controller key) to fetch the results
	function Index(){
		var url = '/exams/staffindex';
		$.get(url,{current:exams.length,parentkey:parentkey,childkey:childkey})
		.done(function(data){
			console.log(data);
			var results = $.parseJSON(data);
				if(results['Result']=="Success"){
					for(var i in results['Exams']){
						var exam = results['Exams'][i];
						exams.push(exam);
						AddRow(exam);
					}
				}else
					alert(results['Error']);
		});
	}

	function DeleteExam(safekey, obj){
		$.ajax({
			type: 'POST',
			url: '/exams/staffdelete',
			data:{
				'safekey':safekey,
				'parentkey':parentkey,
				'childkey':childkey
			},
			success: function(data){
				var results = $.parseJSON(data);
				if(results['Result']=="Success"){
					obj.closest('tr').remove();
				}else
					alert(results['Error']);
			},
			error: function(){
				alert("No connection");
			}
		});
	}

	function AddRow(row){
		table.append('<tr><td><a href="/exams/view/'+row.safekey+'">'+row.name+'</a></td><td><a href="/exams/staffedit/'+row.safekey+'">Edit</a></td><td><a href="#" class="deleteExam" data-id="'+row.safekey+'">Delete</a></td></tr>');
	}
}