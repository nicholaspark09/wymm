var Alerts = new function AlertsManager(){
	
	var that = this;
	var alerts = [];
	var modal;
	var alertButton = $(document).find("#alertClick");
	var moreButton;
	var errorLabel;
	var notiDisplay;
	var table;

	this.Init = function(table_id){
		setupModal();
		GetCount();
		table = $(document).find("#alerts_table");
		alertButton.click(function(e){
			e.preventDefault();
			Index();
		});

		$(document).on("click",".deleteAlert", function(e){
			e.preventDefault();
			var $that = $(this);
			var key = $(this).attr('data-id');
			$that.attr('disabled',true);
			$that.text("Deleting...");
			var url = "/alerts/delete/"+key;
			$.get(url)
			.done(function(data){
				console.log(data);
					var results = $.parseJSON(data);
					if(results['Result']=="Success"){
						$that.closest('tr').remove();
					}else{
						alert(results['Error']);
						$that.attr('disabled',false);
						$that.text("Delete");
					}
			});
		});
	}

	function setupModal(){
		var html = '<div class="modal fade" id="alertsModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">\
      <div class="modal-dialog" role="document">\
        <div class="modal-content">\
          <div class="modal-header">\
        <h4 class="modal-title" id="myNotiLabel">Notifications</h4>\
      </div>\
      <div class="modal-body">\
      <table class="table table-hover"style="overflow:auto;">\
      	<tbody id="alerts_table" style="overflow-y:auto;width:100%;">\
      	</tbody>\
      	</table>\
      </div>\
      <div class="modal-footer">\
      	<div id="alerterrorLabel"></div><br />\
        <a href="#" type="button" class="btn btn-warning" id="viewMoreAlerts">More</a>\
        <button type="button" class="btn btn-default" data-dismiss="modal" id="closeAlerts">Close</button>\
      </div>\
      </div>\
    </div>\
    </div>';
      $(html).appendTo('body');
      modal = $(document).find("#alertsModal");
      moreButton = $(document).find("#viewMoreAlerts");
      errorLabel = $(document).find("#alerterrorLabel");
      alertButton.closest('li').append('<div style="width:20px;height:20px;border-radius:10px;background:red;color:white;font-weight:bold;float:right;position:absolute;right:-5px;top:10px;padding:2px;" id="notiDisplay"></div>');
      notiDisplay = $(document).find("#notiDisplay");
      
	}

	function Index(){
		showMessage("Loading...",false);
		moreButton.attr('disabled',true);
		$.get('/alerts/index',{current:alerts.length,order:1})
		.done(function(data){
			var results = $.parseJSON(data);
			if(results['Result']=="Success"){
				showMessage("",true);
				if(results['Alerts'].length==0){
					showMessage("No More Notifications",true);
				}else{
					for(var i in results['Alerts']){
						var alert = results['Alerts'][i];
						AddRow(alert);
					}
				}
			}else{
				showMessage(results['Error'],true);
			}
			moreButton.attr('disabled',false);
		});
	}

	function GetCount(){
		$.get('/alerts/index',{current:alerts.length,order:2})
		.done(function(data){
			var results = $.parseJSON(data);
			if(results['Result']=="Success"){
				showMessage('',true);
				notiDisplay.text(results['Count']);
			}else{	
				showMessage(results['Error'],true);
			}
		});
	}

	function AddRow(row){
		if(row.seen ==0)
		{
		table.append('<tr style="background:#1abc9c;color:white;"><td>'+row.name+'<br />'+row.body+'</td><td><a href="#" class="deleteAlert" style="color:#c0392b;" data-id="'+row.safekey+'">Delete</a></td></tr>');			
		}else{
					table.append('<tr style="background:white;color:black;"><td>'+row.name+'<br />'+row.body+'</td><td><a href="#" class="deleteAlert" style="color:#c0392b;" data-id="'+row.safekey+'">Delete</a></td></tr>');
		}	

	}

	function showMessage(text,isError){
		if(isError){
			errorLabel.css({'color':'red'});
		}else
			errorLabel.css({'color':'green'});
		errorLabel.text(text);
	}

	that.Init();
}