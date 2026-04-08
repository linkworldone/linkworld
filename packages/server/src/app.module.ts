import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { ScheduleModule } from "@nestjs/schedule";
import { AuthModule } from "./modules/auth/auth.module";
import { UserModule } from "./modules/user/user.module";
import { OperatorModule } from "./modules/operator/operator.module";
import { BillingModule } from "./modules/billing/billing.module";
import { OracleModule } from "./modules/oracle/oracle.module";
import { NotificationModule } from "./modules/notification/notification.module";

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    ScheduleModule.forRoot(),
    AuthModule,
    UserModule,
    OperatorModule,
    BillingModule,
    OracleModule,
    NotificationModule,
  ],
})
export class AppModule {}
