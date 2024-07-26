use rquickjs::{
    embed, loader::Bundle, CatchResultExt, CaughtError, CaughtResult, Context, Module, Runtime,
};

/// load the `my_module.js` file and name it myModule
static BUNDLE: Bundle = embed! {
    "myModule": "script_module.js",
};

fn main() {
    let rt = Runtime::new().unwrap();
    let ctx = Context::full(&rt).unwrap();

    rt.set_loader(BUNDLE, BUNDLE);
    ctx.with(|ctx| {
        match Module::evaluate(
            ctx.clone(),
            "testModule.js",
            r#"
            // await 123;
            // throw new Error("Some error");
            import { f, foo } from 'myModule';
            // import * as a from 'myModule';
            // console.log(f());
            // if(foo() !== 2){
            //     throw new Error("Function didn't return the correct value");
            // }
        "#,
        ) {
            Ok(p) => match p.finish::<()>().catch(&ctx) {
                Ok(_) => {
                    println!("Module loaded successfully");
                }
                Err(e) => match e {
                    CaughtError::Error(e) => {
                        eprintln!("Error: {}", e.to_string());
                    }
                    CaughtError::Exception(e) => {
                        eprintln!("Exception: {}", e);
                    }
                    CaughtError::Value(e) => {
                        eprintln!("Value: {:?}", e);
                    }
                },
            },
            Err(e) => {
                eprintln!("Error: {}", e);
            }
        }
    })
}
