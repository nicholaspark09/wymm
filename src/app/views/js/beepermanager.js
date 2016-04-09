function BeeperManager(){
	var that = this;
	var beepers = [];
	var controller = "";
	var action = "";
	var table;
	var moreButton;
	var safekey = "";
	var types ={1:"One Time Only",2:"Daily",3:"Weekly",4:"Monthly",5:"Quarterly",6:"Special"};
	var intenttypes = {1:"Students",2:"Parents",3:"Teachers",4:"Consultants",5:"Desk Managers",6:"Editors",7:"Ad",8:"Manager",9:"Everyone"};

	this.Init = function(cont,act,table_id,parent){
		safekey = parent;
		controller = cont;
		action = act;
		table = $(document).find("#"+table_id);
		moreButton = table.closest('table').find(".viewMoreBeepers");
		moreButton.click(function(e){
			e.preventDefault();
			Index();
		});
				Index();
		jQuery('#beeperDue').datetimepicker({format:'Y-m-d H:i'});
		var intended = $(document).find("#beeperIntended");
		for(var i in intenttypes){
			intended.append('<option value="'+i+'">'+intenttypes[i]+'</option>');
		}
		$(document).on("click","#createBeeper", function(e){
			var $that = $(this);
			$that.attr('disabled',true);
			$that.text("Saving...");
			var name = $(document).find("#beeperName").val();
			var type = $(document).find("#beeperType").val();
			var due = $(document).find("#beeperDue").val();
			var body = $(document).find("#beeperBody").val();
			var intended = $(document).find("#beeperIntended").val();
			$.ajax({
				type: "POST",
				url: '/beepers/staffadd',
				data:{
					'parentkey':safekey,
					'name':name,
					'body':body,
					'due':due,
					'type':type,
					'controller':controller,
					'action':action,
					'intended':intended
				},
				success: function(data){
					console.log(data);
					var results = $.parseJSON(data);
					if(results['Result']=="Success"){
						var beeper = results['Beeper'];
						beepers.push(beeper);
						AddRow(beeper);
					}else{
						alert(results['Error']);
						console.log(data);
					}
					$that.attr('disabled',false);
					$that.text("Create");
				},error: function(){
					alert(results['Error']);
					$that.attr('disabled',false);
					$that.text("Create");
				}
			})
		});
	}

	function Index(){
		moreButton.text("Searching...");
		moreButton.attr('disabled',true);
		console.log(safekey);
		$.get('/beepers/staffindex',{current:beepers.length,controller:controller,action:action,parentkey:safekey})
		.done(function(data){
			console.log(data);
			var results = $.parseJSON(data);
			if(results['Result']=="Success"){
				for(var i in results['Beepers']){
					var beep = results['Beepers'][i];
					beepers.push(beep);
					AddRow(beep);
				}
				if(results['Beepers'].length<10){
					moreButton.text("No More");
				}else{
					moreButton.attr('disabled',false);
					moreButton.text("More");
				}
			}else
			{
				alert(results['Error']);
				moreButton.attr('disabled',false);
				moreButton.text("More");
			}
		});
	}

	function AddRow(row){
		table.append('<tr><td><a href="/beepers/staffview/'+row.safekey+'">'+row.name+'</a></td><td>'+row.body+'</td><td>'+types[row.type]+'</td><td>'+row.due.substring(0,16)+'</td><td><a href="#" class="editBeeper" data-id="'+row.safekey+'">Edit</a></td><td><a href="#" class="deleteBeeper" data-id="'+row.safekey+'">Delete</a></td></tr>');
	}

}