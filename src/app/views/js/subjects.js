function SubjectsManager(){
	
	var that = this;
	var subjects = [];
	var locale;
	var callFunction;

	this.Init = function(language,callback){
		locale = language;
		callFunction = callback;
		Index();
	}

	function Index(){
		$.get("/subjects/mobileindex",{current:subjects.length, locale:locale})
			.done(function(data){

				var results = $.parseJSON(data);
				if(results['Result'] == "Success")
				{
					for(var i in results['Subjects'])
					{
						var subject = results['Subjects'][i];
						subjects.push(subject);

					}
					callFunction(subjects);
				}

		});
	}


	this.getSubjects = function(){
		return subjects;
	}

}