function MediaMaker(){
          var that = this;
          var safekey = "";
          var mediaFile = $(document).find("#mediaFile");
          var mediaBody = $(document).find("#mediaBody");
          var mediaName = $(document).find("#mediaName");
          var mediaDescription = $(document).find("#mediaDescription");
          var mediaType = $(document).find("#mediaType");
          var mediaTries = $(document).find("#mediaTries");
          var canvas;
          var FR;
          var mediaArea = $(document).find("#mediaArea");
          var body = "";
          var table = $(document).find("#media_table");
          var controller = "";
          var parentCanvas = $(document).find("#parentCanvas").get(0);
          var audioArea = $(document).find("#audioArea");
  

          this.Init = function(){

            $(document).on("click",".deleteMedia", function(e){
                e.preventDefault();
                var key = $(this).attr('data-id');
                var $that =$(this);
                if(confirm("Delete?")){
                  $that.attr('disabled',true);
                  $that.text("Deleting...");
                  var url = '/multimedia/staffdelete/'+key;
                  $.get(url)
                  .done(function(data){
                    var results = $.parseJSON(data);
                    if(results['Result']=="Success"){
                      $that.closest('tr').remove();
                    }else{
                      alert(results['Error']);
                      $that.attr('disabled',false);
                       $that.text("Delete");
                    }
                  });
                }
            });

            $(document).on("click",".clickUse", function(e){
              e.preventDefault();
              //var url = Recorder.getURL();

              //console.log("The media body is "+mediaBody.val());
            });

            $(document).on("click",".viewAddMedia", function(e){
              e.preventDefault();
              mediaBody.val("");
              mediaBody.hide();
              mediaName.val("");
              mediaDescription.val("");
              mediaFile.hide();
              audioArea.hide();
              controller = $(this).attr('data-controller');
              $(document).find("#multiAdd").modal('show');
              safekey = $(this).attr('data-key');
            });

            $(document).on("change","#mediaType", function(e){
              e.preventDefault();
              mediaArea.html("");
              mediaArea.hide();
              var val = mediaType.val();
              if(val==1){
                mediaBody.attr('placeholder',"Text");
                mediaBody.show();
              }else if(val==2){
                mediaBody.attr('placeholder',"HTML");
                mediaBody.show();
              }else if(val==3 || val==4){
                $(document).find("#parentCanvas").show();
                mediaFile.show();
                mediaBody.hide();
              }else if(val==5){
                mediaFile.hide();
                mediaBody.attr('placeholder',"Youtube ID");
                mediaBody.show();
              }
              else if(val==6){
                mediaBody.hide();
                
              }else if(val==7){
                //Show the audio
                audioArea.show();
              }
            });

            $(document).on("change","#mediaBody", function(e){
              e.preventDefault();
              if(mediaType.val()==1){
                //Reading
                
                mediaArea.html("<p>"+mediaBody.val()+"</p>");
              }else if(mediaType.val()==2){
                
                mediaArea.html(mediaBody.val());
              }else if (mediaType.val()==5){
                
                mediaArea.html('<iframe width="560" height="315" src="https://www.youtube.com/embed/'+mediaBody.val()+'" frameborder="0" allowfullscreen></iframe>');
              }
              body = mediaBody.val();
              mediaArea.show();
            });

            $(document).on("change","#mediaFile", function(e){
              e.preventDefault();
              readFile(this);
            });

            $(document).on("click",".viewMedia", function(e){
              e.preventDefault();
              var id = $(this).attr('data-id');
              var body = $(this).attr('data-body');
              var type = $(this).attr('data-type');
              var view = $(document).find("#mediaView");

              if(type=="1"){
                view.html('<p>'+body+'</p>');
              }else if(type=="2"){
                view.html(body);
              }else if(type=="3"){
                view.html('<img src="data:image/png;base64'+body+'" />');
              }else if(type=="4"){
                
                view.html('<audio controls><source src="'+body+'" type="audio/mp3"></audio>');
              }else if(type=="5"){
                view.html('<iframe width="560" height="315" src="https://www.youtube.com/embed/'+body+'" frameborder="0" allowfullscreen></iframe>');
              }else if(type=="7"){
                console.log(body);
                view.html('<audio controls><source src="'+body+'"></audio>');
              }
              $(document).find("#multiView").modal('show');
           
            });
          }

          function readFile(input){
            if ( input.files && input.files[0] ) {
              if(FR==null)
                FR = new FileReader();
                FR.onload = function(e) {

                  if(mediaType.val()==3){
                    //Image
                    
                    var img = document.createElement("img");
                    img.src = e.target.result;
                    data = e.target.result;
                    img.onload = function(){

                      var MAX_WIDTH = 600;
                      var MAX_HEIGHT= 500;
                      var ctx = parentCanvas.getContext('2d');
                            ctx.drawImage(img,0,0);

                        var dimens = calculateSize(this,MAX_WIDTH,MAX_HEIGHT);
                      parentCanvas.width = dimens.x;
                      parentCanvas.height = dimens.y;
                      var ctx = parentCanvas.getContext("2d");
                      ctx.drawImage(this, 0, 0,dimens.x,dimens.y);
                      body = parentCanvas.toDataURL();
                    //mediaArea.html('<img src="'+e.target.result+'" />');
                    }
                  }else if (mediaType.val()==4){
                    //Audio
                    mediaArea.html('<audio controls><source src="'+e.target.result+'" type="audio/mp3"></audio>');
                    body = e.target.result;
                  }
                  mediaArea.show();
                  
                  
                };       
              FR.readAsDataURL( input.files[0] );
              
            }
          }

          function calculateSize(image, max_width, max_height)
            { 
                var width = image.width;
                      var height = image.height;  
                      if (width > height) {
                          if (width > max_width) {
                                height *= max_width / width;
                                width = max_width;
                          }
                      } else {
                          if (height >max_height) {
                              width *= max_height / height;
                              height = max_height;
                          }
                      }
                      return {x:width, y:height};
            }

          //Save it!
          $(document).on("click","#createMedia", function(e){
            e.preventDefault();
            var $that = $(this);
            if(mediaType.val()==0)
              alert("Please enter in a type");
            else{
              
              
              if(mediaType.val()==7)
                body = mediaBody.val();
              console.log(body);
              $that.text("Saving...");
              $that.attr('disabled',true);
              $that.attr('class','btn btn-warning');
              console.log(body);
              $.ajax({
                type: 'POST',
                url: '/multimedia/staffadd',
                data:{
                  'name':mediaName.val(),
                  'description':mediaDescription.val(),
                  'type':mediaType.val(),
                  'body':body,
                  'tries':mediaTries.val(),
                  'controller':controller,
                  'safekey':safekey
                },
                success: function(data){
                  console.log(data);
                  var results = $.parseJSON(data);
                  if(results['Result']=="Success"){
                    var media = {name:mediaName.val(),description:mediaDescription.val(),type:mediaType.val(),prettybody:body,safekey:results['Key'],tries:mediaTries.val(),controller:'sections'};
                    table.append('<tr><td><a href="#" class="viewMedia" data-type="'+media.type+'" data-body="'+media.prettybody+'" data-id="'+media.safekey+'">'+media.name+'</a></td><td>'+media.description+'</td><td>'+media.tries+'</td><td>'+media.type+'</td><td><a href="/sections/staffedit/'+media.safekey+'">Edit</a></td><td><a href="#" class="deleteMedia" data-id="'+media.safekey+'">Delete</a></td></tr>');


                    if(controller=="problems"){
                        $.ajax({
                          type: 'POST',
                          url: '/problems/updatemulti',
                          data:{
                            'safekey':safekey
                          },
                          success: function(data){
                              console.log(data);
                            var results = $.parseJSON(data);
                            if(results['Result']=="Success"){
                              console.log("Saved");
                            }else
                              alert(results['Error']);
                          },
                          error: function(){
                            alert("No connection");
                          }
                        }); 
                    }

                  }else
                    alert(results['Error']);

                    $(document).find("#multiAdd").modal('hide');
                    $that.text("Save");
                    $that.attr('disabled',false);
                    $that.attr('class','btn btn-info');
                },  
                error: function(){
                  alert("No connection");
                  $that.text("Save");
                    $that.attr('disabled',false);
                    $that.attr('class','btn btn-info');
                }
              });
            }
          });


          this.NewVision = function(cont,safe){
            controller = cont;
            safekey = safe;
            mediaBody.val("");
              mediaBody.hide();
              mediaName.val("");
              mediaDescription.val("");
              mediaFile.hide();
              audioArea.hide();
              $(document).find("#multiAdd").modal('show');
          } 
}
