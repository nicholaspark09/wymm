function CategoryManager(){
			var that = this;
			var categories = [];
			var controller = "";
			var picked_key = "";
			var picked_name = "";
			var select = $(document).find("#selectCategory");

			this.Init = function(pickedkey,pickedname, cont){
				picked_key = pickedkey;
				picked_name = pickedname;
				controller = cont;
				index();
			}

			function index(){
				$.get( "/categories/mobileindex", { current: categories.length, controller: controller } )
  									.done(function( data ) {

  											var results= $.parseJSON(data);
  											if(results['Result']=="Success")
  											{
  												for(var i in results['categories'])
  												{
  													var category = results['categories'][i];
  													categories.push(category);
  													select.append('<option value="'+category.safekey+'">'+category.name+'</option>');
  												}
  											}
  					});
			}

			this.getPicked = function(){
				return picked_key;
			}
		}