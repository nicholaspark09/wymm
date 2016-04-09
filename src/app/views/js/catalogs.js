function CatalogManager(){
	
	var that = this;
	var catalog = [];
	var query;
	var subject = "";
	var table = $(document).find("#catalog_table");
	var moreButton = $(document).find("#moreCatalog");


	this.Init = function(){
		query = $(document).find("#catalogSearch");
		$(document).on("click","#searchCatalog",function(e){
			e.preventDefault();
			var $that = $(this);
			table.html('');
			catalog = [];
			Index();
		});
		Index();

		moreButton.click(function(e){
			e.preventDefault();
			Index();
		});
	}

	function Index(){
		moreButton.attr('disabled',true);
		moreButton.text("Searching");
		$.get('/courses/adminindex',{current:catalog.length,query:query.val(),subject:subject})
		.done(function(data){
			var results = $.parseJSON(data);
			if(results['Result']=="Success")
			{
				for(var i in results['Courses'])
				{
					var course = results['Courses'][i];
					catalog.push(course);
					AddRow(course);
				}
			}else
				alert(results['Error']);
			moreButton.attr('disabled',false);
			moreButton.text("More");
		});

	}

	function AddRow(course){
		table.append('<tr><td><a href="/courses/view/'+course.safekey+'">'+course.name+'</a></td><td><a href="#" class="addCatalog" data-id="'+course.safekey+'">Add</a></td></tr>');
	}



}