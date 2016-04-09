function FixerUpper(){
	
	var that = this;
	var languages = ["en","ko","es","zh","ja"];
	var setLanguage = "en";

	this.Init = function(lan){
		if(lan == ""){
			return false;
		}
		console.log(lan);
		setLanguage = lan;
		languages["en"] = {"Subjects":"Subjects","Campus":"Campus","Books":"Books","Profile":"Profile","Email":"Email","Course":"Course","Login":"Login","Name":"Name","Phone":"Phone","Save":"Save","Cancel":"Cancel","Student":"Student","Parent":"Parent","Pin":"Pin","Language":"Language","Logout":"Logout","Alerts":"Alerts","Messages":"Messages","Register":"Register","Drop":"Drop"};
		languages["ko"] = {"Subjects":"과목","Campus":"캠퍼스","Books":"교재","Profile":"프로필","Email":"이메일","Course":"수업","Login":"로그인","Name":"이름","Phone":"전화번호","Save":"동의","Cancel":"취소","Student":"학생","Parent":"학부모","Pin":"비밀번호","Language":"언어","Logout":"로그아웃","Alerts":"알림","Messages":"메시지","Register":"가입하기","Drop":"삭제하기"};
		languages["es"] = {"Subjects":"Sujetos","Campus":"Campus","Books":"Libros","Profile":"Profile","Email":"Email","Course":"Cursos","Login":"Login","Name":"Nombre","Phone":"Teléfono","Save":"Save","Cancel":"Cancelar","Student":"Estudiante","Parent":"Pariente","Pin":"Pin","Language":"Idioma","Logout":"Logout","Alerts":"Alertas","Messages":"Mensajes","Register":"Registrar","Drop":"Cancelar"};
		$(document).find(".newletter").each(function(e){

			var $that = $(this);
			var $text =$that.attr('data-word');
			if(languages[setLanguage][$text]!==undefined)
			{
				$that.text(languages[setLanguage][$text]);
			}
		});
	}

	this.reRun = function(){
		that.Init(setLanguage);
	}

}