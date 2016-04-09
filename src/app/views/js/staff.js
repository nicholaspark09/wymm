function StaffController(){

			var that = this;
			var staff = [];
			var safekey;
			var name;
			var email;
			var phone;
			var staffRole;
			var modal;
			var table = $(document).find("#roles_table");
			var positions = {1:"Student",2:"Parent",3:"Teacher",4:"Parent Consultant",5:"Desk Manager",6:"Editor",7:"Ad",8:"Manager"}



			this.Init = function(key, addOn){

				safekey = key;

				if(addOn){
					TurnOn();
				}
				Index();

				$(document).on("click",".deleteRole", function(e){
					e.preventDefault();
					var $that = $(this);
					console.log("Did you click it?");
					DeleteRow($that);

				});
			}

			function TurnOn(){

				name = $(document).find("#staffName");
				email = $(document).find("#staffEmail");
				phone = $(document).find("#staffPhone");
				staffRole = $(document).find("#staffRole");
				modal = $(document).find("#myModal");
				controller = $(document).find("#staffController");

				$(document).on("click","#addStaff", function(e){
					e.preventDefault();
					var $that =$(this);
					$that.text("Saving...");
					$that.css({'color':'green'});
					//$that.attr('disabled',true);
					$.ajax({
						type: 'POST',
						url: '/roles/adminadd',
						data:
						{	
							'name':name.val(),
							'email':email.val(),
							'phone':phone.val(),
							'level':staffRole.val(),
							'safekey':safekey,
							'controller':controller.val()
						},
						success: function(data){
							console.log(data);
							var results = $.parseJSON(data);
							if(results['Result']=="Success")
							{
								name.val('');
								email.val('');
								phone.val('');
								modal.modal('hide');
							}
							$that.attr('disabled',false);
							Index();
						},
						error: function(){
							alert("No connection");
							$that.attr('disabled',false);
						}
					});
				});
			}

			function DeleteRow($button){
				console.log("Operating...");
				if(!confirm("Are you sure?")){
					return;
				}
				var id = $button.attr('data-id');
				var url = '/roles/delete/'+id;
				$button.text("Deleting...");
				$button.attr('disabled',true);
				$.get(url)
				.done(function(data){
					var results = $.parseJSON(data);
					if(results['Result']=="Success")
					{
						$button.closest('tr').remove();
						for(var i in staff){
							if(staff[i].safekey == id){
								staff.splice(i,1);
								break;
							}
						}
					}else{
						alert(results['Error']);
						$button.text("Delete");
						$button.attr('disabled',false);
					}
				});
			}

			function Index(){
				$.get('/roles/adminindex',{safekey:safekey,current: staff.length})
				.done(function(data){
					console.log(data);
					var results = $.parseJSON(data);
					if(results['Result']=="Success")
					{
						for(var i in results['Roles'])
						{
							var employee = results['Roles'][i];
							staff.push(employee);
							AddRow(employee);
						}
						
					}else{
						alert(results['Error']);
					}
				});
			}

			function AddRow(role){
				table.append('<tr><td><a href="/roles/view/'+role.safekey+'">'+role.name+'</a></td><td>'+positions[role.level]+'</td><td>'+role.phone+'</td><td><a href="#" class="editRole" data-id="'+role.safekey+'">Edit</a></td><td><a href="#" class="deleteRole" data-id="'+role.safekey+'">Delete</a></td><td><a href="/nfcs/lookupstudent/'+role.safekey+'" class="viewNFC" >View NFC</a></td></tr>');
			}


		}